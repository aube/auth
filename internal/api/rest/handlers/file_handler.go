package handlers

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/aube/auth/internal/api/rest/dto"
	appFile "github.com/aube/auth/internal/application/file"
	"github.com/aube/auth/internal/utils/logger"
	"github.com/rs/zerolog"

	"github.com/gin-gonic/gin"
)

// FileHandler обрабатывает HTTP запросы для работы с файлами
type FileHandler struct {
	service *appFile.FileService
	log     zerolog.Logger
}

// NewFileHandler создает новый экземпляр FileHandler
func NewFileHandler(service *appFile.FileService) *FileHandler {
	return &FileHandler{
		service: service,
		log:     logger.Get().With().Str("handlers", "file_handler").Logger(),
	}
}

// UploadFile обрабатывает загрузку файла
func (h *FileHandler) UploadFile(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		h.log.Debug().Err(err).Msg("UploadFile1")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		h.log.Debug().Err(err).Msg("UploadFile2")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}
	defer file.Close()

	uploadedFile, err := h.service.Upload(
		c.Request.Context(),
		fileHeader.Filename,
		fileHeader.Header.Get("Content-Type"),
		fileHeader.Size,
		file,
	)
	if err != nil {
		h.log.Debug().Err(err).Msg("UploadFile3")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	c.JSON(http.StatusCreated, dto.NewFileResponse(uploadedFile))
}

// DownloadFile обрабатывает скачивание файла
func (h *FileHandler) DownloadFile(c *gin.Context) {
	fileID := c.Query("id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File ID is required"})
		return
	}

	file, err := h.service.GetByID(c.Request.Context(), fileID)
	if err != nil {
		h.log.Debug().Err(err).Msg("DownloadFile1")
		if errors.Is(err, appFile.ErrFileNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file"})
		return
	}

	content, err := h.service.Download(c.Request.Context(), file)
	if err != nil {
		h.log.Debug().Err(err).Msg("DownloadFile2")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}
	defer content.Close()

	c.Header("Content-Disposition", "attachment; filename="+file.Name)
	c.Header("Content-Type", file.ContentType)
	c.Header("Content-Length", strconv.FormatInt(file.Size, 10))

	c.Stream(func(w io.Writer) bool {
		if _, err := io.Copy(w, content); err != nil {
			return false
		}
		return false
	})
}

// ListFiles возвращает список всех файлов
func (h *FileHandler) ListFiles(c *gin.Context) {
	files, err := h.service.List(c.Request.Context())
	if err != nil {
		h.log.Debug().Err(err).Msg("ListFiles")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list files"})
		return
	}

	response := make([]dto.FileResponse, len(files))
	for i, file := range files {
		response[i] = dto.NewFileResponse(file)
	}

	c.JSON(http.StatusOK, response)
}

// DeleteFile обрабатывает удаление файла
func (h *FileHandler) DeleteFile(c *gin.Context) {
	fileID := c.Query("id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File ID is required"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), fileID); err != nil {
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
