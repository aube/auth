// Package file provides file storage repository interfaces and implementations.
package file

import (
	"context"
	"errors"
	"io"

	"github.com/aube/auth/internal/domain/entities"
)

// ErrFileNotFound is returned when a requested file cannot be located in storage.
var ErrFileNotFound = errors.New("file not found")

// FileRepository defines the interface for file persistence operations.
// Implementations should handle actual file storage (e.g., disk, cloud storage).
//
// Methods:
//
//   - Save: Stores a new file with metadata and content
//     ctx: Context for request cancellation/timeout
//     file: File metadata entity
//     data: File content stream
//     Returns: error on failure
//
//   - FindAll: Retrieves metadata for all stored files
//     ctx: Context for request cancellation/timeout
//     Returns: (*entities.Files, error)
//
//   - Delete: Removes a file by its identifier
//     ctx: Context for request cancellation/timeout
//     id: File UUID
//     Returns: error on failure
//
//   - GetFileContent: Retrieves file content as a stream
//     ctx: Context for request cancellation/timeout
//     uuid: File identifier
//     Returns: (io.ReadCloser, error) - caller must close the stream
type FileRepository interface {
	Save(ctx context.Context, file *entities.File, data io.Reader) error
	FindAll(ctx context.Context) (*entities.Files, error)
	Delete(ctx context.Context, id string) error
	GetFileContent(ctx context.Context, uuid string) (io.ReadCloser, error)
}
