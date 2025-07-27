package upload_test

import (
	"context"

	"github.com/aube/auth/internal/application/dto"
	"github.com/aube/auth/internal/domain/entities"
	"github.com/stretchr/testify/mock"
)

type UploadRepository struct {
	mock.Mock
}

func (m *UploadRepository) Create(ctx context.Context, userID int64, upload *entities.Upload) error {
	args := m.Called(ctx, userID, upload)
	return args.Error(0)
}

func (m *UploadRepository) ListByUserID(ctx context.Context, userID int64, offset, limit int) (*entities.Uploads, *dto.Pagination, error) {
	args := m.Called(ctx, userID, offset, limit)
	return args.Get(0).(*entities.Uploads), args.Get(1).(*dto.Pagination), args.Error(2)
}

func (m *UploadRepository) GetByUUID(ctx context.Context, uuid string, userID int64) (*entities.Upload, error) {
	args := m.Called(ctx, uuid, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Upload), args.Error(1)
}

func (m *UploadRepository) GetByName(ctx context.Context, name string, userID int64) (*entities.Upload, error) {
	args := m.Called(ctx, name, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Upload), args.Error(1)
}

func (m *UploadRepository) Delete(ctx context.Context, uuid string, userID int64) error {
	return m.Called(ctx, uuid, userID).Error(0)
}

func (m *UploadRepository) DeleteForce(ctx context.Context, uuid string, userID int64) error {
	return m.Called(ctx, uuid, userID).Error(0)
}
