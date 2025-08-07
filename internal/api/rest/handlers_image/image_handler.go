package handlers_image

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/aube/auth/internal/application/dto"
	appFile "github.com/aube/auth/internal/application/file"
	appImage "github.com/aube/auth/internal/application/upload"
	"github.com/aube/auth/internal/domain/entities"
	"github.com/aube/auth/internal/utils/logger"
	"github.com/rs/zerolog"

	"github.com/gin-gonic/gin"
)

type FileService interface {
	Delete(ctx context.Context, id string) error
	Download(ctx context.Context, uuid string) (io.ReadCloser, error)
	Upload(ctx context.Context, size int64, data io.Reader) (*entities.File, error)
}

type ImageService interface {
	Delete(ctx context.Context, uuid string, userID int64) error
	DeleteForce(ctx context.Context, uuid string, userID int64) error
	GetByName(ctx context.Context, name string, userID int64) (*entities.Image, error)
	GetByUUID(ctx context.Context, uuid string, userID int64) (*entities.Image, error)
	ListByUserID(ctx context.Context, userID int64, offset int, limit int, params map[string]any) (*entities.Images, *dto.Pagination, error)
	RegisterUploadedImage(ctx context.Context, userID int64, file *entities.File, name string, category string, contentType string, description string) (*entities.Image, error)
}

type ImageHandler interface {
	DeleteFile(c *gin.Context)
	DownloadFile(c *gin.Context)
	ListFiles(c *gin.Context)
	ImageFile(c *gin.Context)
}

// SavedFile implements structure for saved file results.
// File: FileService.Upload operation result.
// Filename: Data from fileHeader.Filename.
// ContentType: Content-Type uploaded file.
type SavedFile struct {
	File        *entities.File
	Filename    string
	ContentType string
}

// Handler implements UploadHandler for handling file-related HTTP requests.
// FileService: Service for file storage operations.
// ImageService: Service for upload metadata operations.
// log: Logger instance for the handler.
type Handler struct {
	FileService  FileService
	ImageService ImageService
	log          zerolog.Logger
}

// NewHandler создает новый экземпляр Handler
func NewImageHandler(FileService FileService, ImageService ImageService) *Handler {
	return &Handler{
		FileService:  FileService,
		ImageService: ImageService,
		log:          logger.Get().With().Str("handlers", "file_handler").Logger(),
	}
}

// UploadImage обрабатывает загрузку файла
func (h *Handler) UploadImage(c *gin.Context) {

	userID := c.GetInt("userID")
	description := c.PostForm("description")
	category := c.PostForm("category")

	savedFile, err := h.saveFile(c, "file", userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}

	upload, err := h.ImageService.RegisterUploadedImage(
		c.Request.Context(),
		int64(userID),
		savedFile.File,
		savedFile.Filename,
		category,
		savedFile.ContentType,
		description,
	)
	if err != nil {
		h.log.Debug().Err(err).Msg("ImageFile5")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write upload file into DB"})
		return
	}

	c.JSON(http.StatusCreated, dto.NewImageResponse(upload))
}

// DownloadFile handles file download requests.
// Supports lookup by UUID or filename and enforces user ownership.
func (h *Handler) DownloadFile(c *gin.Context) {

	upload, err := h.getUpload(c)
	if err != nil {
		return
	}

	content, err := h.FileService.Download(c.Request.Context(), upload.UUID)
	if err != nil {
		if errors.Is(err, appFile.ErrFileNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found on FS"})
			return
		}
		h.log.Debug().Err(err).Msg("DownloadFile2")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}
	defer content.Close()

	c.Header("Content-Disposition", "attachment; filename="+upload.Name)
	c.Header("Content-Type", upload.ContentType)
	c.Header("Content-Length", strconv.FormatInt(upload.Size, 10))

	c.Stream(func(w io.Writer) bool {
		if _, err := io.Copy(w, content); err != nil {
			return false
		}
		return false
	})
}

// ListFiles retrieves a paginated list of files uploaded by the user.
// Uses PaginationMiddleware for offset/limit handling.
func (h *Handler) ListFiles(c *gin.Context) {

	userID := c.GetInt("userID")
	offset := c.GetInt("offset")
	limit := c.GetInt("limit")

	params := make(map[string]any)
	if c.Query("updated_at") != "" {
		params["updated_at >="] = c.Query("updated_at")
	}
	params["deleted"] = "false"

	uploads, pagination, err := h.ImageService.ListByUserID(c.Request.Context(), int64(userID), offset, limit, params)
	if err != nil {
		h.log.Debug().Err(err).Msg("ListFiles")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list files"})
		return
	}

	rows := make([]dto.ImageResponse, len(*uploads))
	for i, upload := range *uploads {
		rows[i] = dto.NewImageResponse(&upload)
	}

	c.JSON(http.StatusOK, gin.H{
		"rows":       rows,
		"pagination": pagination,
	})
}

func (h *Handler) DeleteFile(c *gin.Context) {

	userID := c.GetInt("userID")
	UUID := c.Query("uuid")

	if UUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File UUID is required"})
		return
	}
	if err := h.ImageService.Delete(c.Request.Context(), UUID, int64(userID)); err != nil {
		h.log.Debug().Err(err).Msg("DeleteFile")
		c.JSON(http.StatusBadRequest, gin.H{"error": "File UUID is can't be deleted"})
		return
	}

	if err := h.FileService.Delete(c.Request.Context(), UUID); err != nil {
		h.log.Debug().Err(err).Msg("DeleteFile")
		if errors.Is(err, appFile.ErrFileNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) cleanupBeforeCreate(ctx context.Context, name string, userID int64) error {

	upload, err := h.ImageService.GetByName(ctx, name, userID)

	if err != nil {
		return nil // file not found
	}

	if err := h.ImageService.DeleteForce(ctx, upload.UUID, userID); err != nil {
		return err
	}

	if err := h.FileService.Delete(ctx, upload.UUID); err != nil {
		return err
	}

	return nil
}

// saveFile write file to FS via FileService.Upload.
func (h *Handler) saveFile(c *gin.Context, fieldName string, userID int) (*SavedFile, error) {

	fileHeader, err := c.FormFile(fieldName)
	if err != nil {
		h.log.Debug().Err(err).Msg("saveFile1")
		return nil, err
	}

	err = h.cleanupBeforeCreate(c.Request.Context(), fileHeader.Filename, int64(userID))
	if err != nil {
		h.log.Debug().Err(err).Msg("saveFile2")
		return nil, err
	}

	uploadingFile, err := fileHeader.Open()
	if err != nil {
		h.log.Debug().Err(err).Msg("saveFile3")
		return nil, err
	}
	defer uploadingFile.Close()

	file, err := h.FileService.Upload(
		c.Request.Context(),
		fileHeader.Size,
		uploadingFile,
	)
	if err != nil {
		h.log.Debug().Err(err).Msg("saveFile4")
		return nil, err
	}

	return &SavedFile{
		file,
		fileHeader.Filename,
		fileHeader.Header.Get("Content-Type"),
	}, nil
}

// getUpload return data from uploads table DB by uuid or filename.
func (h *Handler) getUpload(c *gin.Context) (*entities.Image, error) {

	userID := c.GetInt("userID")
	UUID := c.Query("uuid")
	name := c.Query("name")

	if UUID == "" && name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File UUID or Name is required"})
		return nil, errors.New("bad request")
	}

	var upload *entities.Image
	var err error

	if UUID != "" {
		upload, err = h.ImageService.GetByUUID(c.Request.Context(), UUID, int64(userID))
	} else {
		upload, err = h.ImageService.GetByName(c.Request.Context(), name, int64(userID))
	}

	if err != nil {
		if errors.Is(err, appImage.ErrFileNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found in DB"})
			return nil, err
		}
		h.log.Debug().Err(err).Msg("DownloadFile1")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file"})
		return nil, err
	}

	return upload, err
}
