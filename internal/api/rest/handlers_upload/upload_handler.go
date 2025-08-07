// Package handlers_upload provides handlers for file upload and management operations.
package handlers_upload

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/aube/auth/internal/application/dto"
	appFile "github.com/aube/auth/internal/application/file"
	appUpload "github.com/aube/auth/internal/application/upload"
	"github.com/aube/auth/internal/domain/entities"
	"github.com/aube/auth/internal/utils/logger"
	"github.com/rs/zerolog"

	"github.com/gin-gonic/gin"
)

// FileService defines the interface for file storage operations (upload, download, delete).
type FileService interface {
	Delete(ctx context.Context, id string) error
	Download(ctx context.Context, uuid string) (io.ReadCloser, error)
	Upload(ctx context.Context, size int64, data io.Reader) (*entities.File, error)
}

// UploadService defines the interface for upload metadata operations (CRUD and listing).
type UploadService interface {
	Delete(ctx context.Context, uuid string, userID int64) error
	DeleteForce(ctx context.Context, uuid string, userID int64) error
	GetByName(ctx context.Context, name string, userID int64) (*entities.Upload, error)
	GetByUUID(ctx context.Context, uuid string, userID int64) (*entities.Upload, error)
	ListByUserID(ctx context.Context, userID int64, offset int, limit int, params map[string]any) (*entities.Uploads, *dto.Pagination, error)
	RegisterUploadedFile(ctx context.Context, userID int64, file *entities.File, name string, category string, contentType string, description string) (*entities.Upload, error)
}

type UploadHandler interface {
	DeleteFile(c *gin.Context)
	DownloadFile(c *gin.Context)
	ListFiles(c *gin.Context)
	UploadFile(c *gin.Context)
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
// UploadService: Service for upload metadata operations.
// log: Logger instance for the handler.
type Handler struct {
	FileService   FileService
	UploadService UploadService
	log           zerolog.Logger
}

// NewHandler создает новый экземпляр Handler
func NewUploadHandler(FileService FileService, UploadService UploadService) *Handler {
	return &Handler{
		FileService:   FileService,
		UploadService: UploadService,
		log:           logger.Get().With().Str("handlers", "file_handler").Logger(),
	}
}

// UploadFile handles file upload requests.
// Validates user authentication, processes the file, and stores metadata.
func (h *Handler) UploadFile(c *gin.Context) {

	userID := c.GetInt("userID")
	description := c.PostForm("description")
	category := c.PostForm("category")

	savedFile, err := h.saveFile(c, "file", userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}

	upload, err := h.UploadService.RegisterUploadedFile(
		c.Request.Context(),
		int64(userID),
		savedFile.File,
		savedFile.Filename,
		category,
		savedFile.ContentType,
		description,
	)
	if err != nil {
		h.log.Debug().Err(err).Msg("UploadFile5")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write upload file into DB"})
		return
	}

	c.JSON(http.StatusCreated, dto.NewUploadResponse(upload))
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

	uploads, pagination, err := h.UploadService.ListByUserID(c.Request.Context(), int64(userID), offset, limit, params)
	if err != nil {
		h.log.Debug().Err(err).Msg("ListFiles")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list files"})
		return
	}

	rows := make([]dto.UploadResponse, len(*uploads))
	for i, upload := range *uploads {
		rows[i] = dto.NewUploadResponse(&upload)
	}

	c.JSON(http.StatusOK, gin.H{
		"rows":       rows,
		"pagination": pagination,
	})
}

// DeleteFile handles file deletion requests.
// Validates user ownership before deleting both file and metadata.
func (h *Handler) DeleteFile(c *gin.Context) {

	userID := c.GetInt("userID")
	UUID := c.Query("uuid")

	if UUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File UUID is required"})
		return
	}
	if err := h.UploadService.Delete(c.Request.Context(), UUID, int64(userID)); err != nil {
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

// cleanupBeforeCreate ensures no duplicate filenames exist for the user by cleaning up existing entries.
func (h *Handler) cleanupBeforeCreate(ctx context.Context, name string, userID int64) error {

	upload, err := h.UploadService.GetByName(ctx, name, userID)

	if err != nil {
		return nil // file not found
	}

	if err := h.UploadService.DeleteForce(ctx, upload.UUID, userID); err != nil {
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
func (h *Handler) getUpload(c *gin.Context) (*entities.Upload, error) {

	userID := c.GetInt("userID")
	UUID := c.Query("uuid")
	name := c.Query("name")

	if UUID == "" && name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File UUID or Name is required"})
		return nil, errors.New("bad request")
	}

	var upload *entities.Upload
	var err error

	if UUID != "" {
		upload, err = h.UploadService.GetByUUID(c.Request.Context(), UUID, int64(userID))
	} else {
		upload, err = h.UploadService.GetByName(c.Request.Context(), name, int64(userID))
	}

	if err != nil {
		if errors.Is(err, appUpload.ErrFileNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found in DB"})
			return nil, err
		}
		h.log.Debug().Err(err).Msg("DownloadFile1")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file"})
		return nil, err
	}

	return upload, err
}
