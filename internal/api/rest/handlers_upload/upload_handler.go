package handlers_upload

import (
	"context"
	"errors"
	"fmt"
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

// UploadHandler обрабатывает HTTP запросы для работы с файлами
type UploadHandler struct {
	FileService   *appFile.FileService
	UploadService *appUpload.UploadService
	log           zerolog.Logger
}

// NewUploadHandler создает новый экземпляр UploadHandler
func NewUploadHandler(FileService *appFile.FileService, UploadService *appUpload.UploadService) *UploadHandler {
	return &UploadHandler{
		FileService:   FileService,
		UploadService: UploadService,
		log:           logger.Get().With().Str("handlers", "file_handler").Logger(),
	}
}

// UploadFile обрабатывает загрузку файла
func (h *UploadHandler) UploadFile(c *gin.Context) {
	uID, exists := c.Get("userID")
	if !exists {
		h.log.Debug().Msg("GetProfile not exists")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, ok := uID.(int64)
	if !ok {
		h.log.Debug().Msg("GetProfile not ok")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized2"})
		return
	}

	description := c.PostForm("description")
	category := c.PostForm("category")

	fileHeader, err := c.FormFile("file")
	if err != nil {
		h.log.Debug().Err(err).Msg("UploadFile1")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}

	err = h.cleanupBeforeCreate(c.Request.Context(), fileHeader.Filename, userID)
	if err != nil {
		h.log.Debug().Err(err).Msg("UploadFile2")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}

	uploadingFile, err := fileHeader.Open()
	if err != nil {
		h.log.Debug().Err(err).Msg("UploadFile3")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}
	defer uploadingFile.Close()

	file, err := h.FileService.Upload(
		c.Request.Context(),
		fileHeader.Size,
		uploadingFile,
	)
	if err != nil {
		h.log.Debug().Err(err).Msg("UploadFile4")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	upload, err := h.UploadService.RegisterUploadedFile(
		c.Request.Context(),
		userID,
		file,
		fileHeader.Filename,
		category,
		fileHeader.Header.Get("Content-Type"),
		description,
	)
	if err != nil {
		h.log.Debug().Err(err).Msg("UploadFile5")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write upload file into DB"})
		return
	}

	c.JSON(http.StatusCreated, dto.NewUploadResponse(upload))
}

func (h *UploadHandler) DownloadFile(c *gin.Context) {
	uID, exists := c.Get("userID")
	if !exists {
		h.log.Debug().Msg("GetProfile not exists")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, ok := uID.(int64)
	if !ok {
		h.log.Debug().Msg("GetProfile not ok")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized2"})
		return
	}

	UUID := c.Query("uuid")
	name := c.Query("name")

	if UUID == "" && name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File UUID or Name is required"})
		return
	}

	var upload *entities.Upload
	var err error

	if UUID != "" {
		upload, err = h.UploadService.GetByUUID(c.Request.Context(), UUID, userID)
	} else {
		upload, err = h.UploadService.GetByName(c.Request.Context(), name, userID)
	}
	fmt.Println(upload)
	if err != nil {
		if errors.Is(err, appUpload.ErrFileNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found in DB"})
			return
		}
		h.log.Debug().Err(err).Msg("DownloadFile1")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file"})
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

func (h *UploadHandler) ListFiles(c *gin.Context) {
	uID, exists := c.Get("userID")
	if !exists {
		h.log.Debug().Msg("GetProfile not exists")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, ok := uID.(int64)
	if !ok {
		h.log.Debug().Msg("GetProfile not ok")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized2"})
		return
	}

	offset := c.GetInt("offset")
	limit := c.GetInt("limit")

	uploads, pagination, err := h.UploadService.ListByUserID(c.Request.Context(), userID, offset, limit)
	if err != nil {
		h.log.Debug().Err(err).Msg("ListFiles")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list files"})
		return
	}

	rows := make([]dto.UploadResponse, len(*uploads))
	for i, upload := range *uploads {
		rows[i] = dto.NewUploadResponse(&upload)
	}

	type resp struct {
		rows       []dto.UploadResponse
		pagination *dto.Pagination
	}

	fmt.Println(rows)
	fmt.Println(pagination)
	fmt.Println(pagination)
	fmt.Println(pagination)
	c.JSON(http.StatusOK, gin.H{
		"rows":       rows,
		"pagination": pagination,
	})
	// c.JSON(http.StatusOK, &resp{
	// 	rows:       rows,
	// 	pagination: pagination,
	// })
}

func (h *UploadHandler) DeleteFile(c *gin.Context) {
	uID, exists := c.Get("userID")
	if !exists {
		h.log.Debug().Msg("GetProfile not exists")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, ok := uID.(int64)
	if !ok {
		h.log.Debug().Msg("GetProfile not ok")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized2"})
		return
	}

	UUID := c.Query("uuid")
	if UUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File UUID is required"})
		return
	}

	if err := h.UploadService.Delete(c.Request.Context(), UUID, userID); err != nil {
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

func (h *UploadHandler) cleanupBeforeCreate(ctx context.Context, name string, userID int64) error {

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
