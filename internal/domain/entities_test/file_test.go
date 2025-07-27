package entities_test

import (
	"testing"

	"github.com/aube/auth/internal/domain/entities"
	"github.com/stretchr/testify/assert"
)

func TestNewFile(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		filePath string
		size     int64
	}{
		{"regular file", "test.txt", "/uploads/test.txt", 1024},
		{"empty name", "", "/uploads/", 0},
		{"large file", "bigfile.iso", "/storage/bigfile.iso", 1024 * 1024 * 1024},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := entities.NewFile(tt.fileName, tt.filePath, tt.size)
			assert.Equal(t, tt.fileName, file.Name)
			assert.Equal(t, tt.filePath, file.Path)
			assert.Equal(t, tt.size, file.Size)
		})
	}
}

func TestFilesSlice(t *testing.T) {
	files := entities.Files{
		{Name: "file1.txt", Size: 100},
		{Name: "file2.txt", Size: 200},
	}

	assert.Len(t, files, 2)
	assert.Equal(t, "file1.txt", files[0].Name)
	assert.Equal(t, int64(200), files[1].Size)
}
