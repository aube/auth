package entities_test

import (
	"testing"
	"time"

	"github.com/aube/auth/internal/domain/entities"
	"github.com/stretchr/testify/assert"
)

func TestNewUpload(t *testing.T) {
	now := time.Now()
	file := &entities.File{
		Name: "file-uuid",
		Path: "/uploads/file-uuid",
		Size: 1024,
	}

	upload := entities.NewUpload(file, 1, 100, "original.txt", "docs", "text/plain", "test file", now)

	assert.Equal(t, int64(1), upload.ID)
	assert.Equal(t, int64(100), upload.UserID)
	assert.Equal(t, "file-uuid", upload.UUID)
	assert.Equal(t, "original.txt", upload.Name)
	assert.Equal(t, "docs", upload.Category)
	assert.Equal(t, int64(1024), upload.Size)
	assert.Equal(t, "text/plain", upload.ContentType)
	assert.Equal(t, "/uploads/file-uuid", upload.Path)
	assert.Equal(t, "test file", upload.Description)
	assert.Equal(t, now, upload.UploadedAt)
}

func TestUploadsSlice(t *testing.T) {
	uploads := entities.Uploads{
		{ID: 1, Name: "file1.txt"},
		{ID: 2, Name: "file2.txt"},
	}

	assert.Len(t, uploads, 2)
	assert.Equal(t, int64(1), uploads[0].ID)
	assert.Equal(t, "file2.txt", uploads[1].Name)
}
