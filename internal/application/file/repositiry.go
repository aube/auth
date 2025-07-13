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
	FindByID(ctx context.Context, id string) (*entities.File, error)
	FindAll(ctx context.Context) ([]*entities.File, error)
	Delete(ctx context.Context, id string) error
	GetFileContent(ctx context.Context, file *entities.File) (io.ReadCloser, error)
}
