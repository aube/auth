package dto_test

import (
	"testing"

	"github.com/aube/auth/internal/application/dto"
	"github.com/aube/auth/internal/domain/entities"
	"github.com/stretchr/testify/assert"
)

func TestNewUploadResponse(t *testing.T) {
	upload := &entities.Upload{
		UUID:        "test-uuid",
		Name:        "test.txt",
		Category:    "docs",
		Size:        1234,
		ContentType: "text/plain",
		Description: "test file",
	}

	resp := dto.NewUploadResponse(upload)

	assert.Equal(t, "test-uuid", resp.UUID)
	assert.Equal(t, "test.txt", resp.Name)
	assert.Equal(t, "docs", resp.Category)
	assert.Equal(t, int64(1234), resp.Size)
	assert.Equal(t, "text/plain", resp.ContentType)
	assert.Equal(t, "test file", resp.Description)
}

func TestUploadResponse_Empty(t *testing.T) {
	resp := dto.UploadResponse{}

	assert.Empty(t, resp.UUID)
	assert.Empty(t, resp.Name)
	assert.Empty(t, resp.Category)
	assert.Equal(t, int64(0), resp.Size)
	assert.Empty(t, resp.ContentType)
	assert.Empty(t, resp.Description)
}
