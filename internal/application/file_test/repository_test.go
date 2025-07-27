package file_test

import (
	"context"
	"io"

	"github.com/aube/auth/internal/domain/entities"
	"github.com/stretchr/testify/mock"
)

type FileRepository struct {
	mock.Mock
}

func (m *FileRepository) Save(ctx context.Context, file *entities.File, data io.Reader) error {
	args := m.Called(ctx, file, data)
	return args.Error(0)
}

func (m *FileRepository) FindAll(ctx context.Context) (*entities.Files, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Files), args.Error(1)
}

func (m *FileRepository) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

func (m *FileRepository) GetFileContent(ctx context.Context, uuid string) (io.ReadCloser, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(io.ReadCloser), args.Error(1)
}
