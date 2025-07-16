package handlers

import (
	"net/http"

	appUpload "github.com/aube/auth/internal/application/upload"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	UploadService *appUpload.UploadService
	jwtSecret     []byte
}

func NewUploadHandler(UploadService *appUpload.UploadService, jwtSecret string) *UploadHandler {
	return &UploadHandler{
		UploadService: UploadService,
		jwtSecret:     []byte(jwtSecret),
	}
}

func (h *UploadHandler) ListUploads(c *gin.Context) {

	ctx := c.Request.Context()
	userID, _ := c.Get("userID")

	uidStr, ok := userID.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userID type"})
		return
	}

	uploads, err := h.UploadService.ListByUserID(ctx, uidStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, uploads)
}
