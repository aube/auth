package file

import (
	"context"
	"errors"
	"io"

	"github.com/aube/auth/internal/domain/entities"
)

var ErrFileNotFound = errors.New("file not found")

type FileRepository interface {
	Save(ctx context.Context, file *entities.File, data io.Reader) error
	FindAll(ctx context.Context) (*entities.Files, error)
	Delete(ctx context.Context, id string) error
	GetFileContent(ctx context.Context, uuid string) (io.ReadCloser, error)
}
