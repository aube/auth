package upload_test

import (
	"context"
	"testing"

	"github.com/aube/auth/internal/application/dto"
	appUpload "github.com/aube/auth/internal/application/upload"
	"github.com/aube/auth/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUploadService_RegisterUploadedFile_Success(t *testing.T) {
	// Setup
	mockRepo := new(UploadRepository)
	service := appUpload.NewUploadService(mockRepo)

	// Test data
	userID := int64(1)
	file := &entities.File{
		Name: "file-uuid",
		Size: 1024,
	}
	// Mock expectations
	mockRepo.On("Create", mock.Anything, userID, mock.AnythingOfType("*entities.Upload")).
		Run(func(args mock.Arguments) {
			uploadArg := args.Get(2).(*entities.Upload)
			// assert.Equal(t, file.Name, uploadArg.File.Name)
			assert.Equal(t, userID, uploadArg.UserID)
			assert.Equal(t, "test.txt", uploadArg.Name)
			assert.Equal(t, "docs", uploadArg.Category)
			assert.Equal(t, "text/plain", uploadArg.ContentType)
			assert.Equal(t, "test file", uploadArg.Description)
		}).
		Return(nil)

	// Execute
	result, err := service.RegisterUploadedFile(
		context.Background(),
		userID,
		file,
		"test.txt",
		"docs",
		"text/plain",
		"test file",
	)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	// assert.Equal(t, file.Name, result.File.Name)
	mockRepo.AssertExpectations(t)
}

func TestUploadService_ListByUserID_Success(t *testing.T) {
	// Setup
	mockRepo := new(UploadRepository)
	service := appUpload.NewUploadService(mockRepo)

	// Test data
	userID := int64(1)
	offset := 0
	limit := 10
	params := make(map[string]any)
	expectedUploads := &entities.Uploads{
		entities.Upload{UUID: "uuid1"},
		entities.Upload{UUID: "uuid2"},
	}
	expectedPagination := &dto.Pagination{Total: 2, Page: 1, Size: 10}

	// Mock expectations
	mockRepo.On("ListByUserID", mock.Anything, userID, offset, limit, params).
		Return(expectedUploads, expectedPagination, nil)

	// Execute
	uploads, pagination, err := service.ListByUserID(context.Background(), userID, offset, limit, params)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedUploads, uploads)
	assert.Equal(t, expectedPagination, pagination)
	mockRepo.AssertExpectations(t)
}

func TestUploadService_GetByUUID_Success(t *testing.T) {
	// Setup
	mockRepo := new(UploadRepository)
	service := appUpload.NewUploadService(mockRepo)

	// Test data
	userID := int64(1)
	uuid := "test-uuid"
	expectedUpload := &entities.Upload{UUID: uuid}

	// Mock expectations
	mockRepo.On("GetByUUID", mock.Anything, uuid, userID).
		Return(expectedUpload, nil)

	// Execute
	upload, err := service.GetByUUID(context.Background(), uuid, userID)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, expectedUpload, upload)
	mockRepo.AssertExpectations(t)
}

func TestUploadService_GetByUUID_NotFound(t *testing.T) {
	// Setup
	mockRepo := new(UploadRepository)
	service := appUpload.NewUploadService(mockRepo)

	// Test data
	userID := int64(1)
	uuid := "not-found"

	// Mock expectations
	mockRepo.On("GetByUUID", mock.Anything, uuid, userID).
		Return(nil, appUpload.ErrFileNotFound)

	// Execute
	upload, err := service.GetByUUID(context.Background(), uuid, userID)

	// Assert
	assert.Nil(t, upload)
	assert.Equal(t, appUpload.ErrFileNotFound, err)
	mockRepo.AssertExpectations(t)
}

func TestUploadService_DeleteForce_Success(t *testing.T) {
	// Setup
	mockRepo := new(UploadRepository)
	service := appUpload.NewUploadService(mockRepo)

	// Test data
	userID := int64(1)
	uuid := "test-uuid"

	// Mock expectations
	mockRepo.On("DeleteForce", mock.Anything, uuid, userID).Return(nil)

	// Execute
	err := service.DeleteForce(context.Background(), uuid, userID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockRepo.AssertNotCalled(t, "GetByUUID") // DeleteForce не должен проверять существование
}

func TestUploadService_CanBeDeleted_Success(t *testing.T) {
	// Setup
	mockRepo := new(UploadRepository)
	service := appUpload.NewUploadService(mockRepo)

	// Test data
	userID := int64(1)
	uuid := "test-uuid"
	expectedUpload := &entities.Upload{UUID: uuid, UserID: userID}

	// Mock expectations
	mockRepo.On("GetByUUID", mock.Anything, uuid, userID).Return(expectedUpload, nil)

	// Execute
	err := service.CanBeDeleted(context.Background(), uuid, userID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUploadService_CanBeDeleted_NotFound(t *testing.T) {
	// Setup
	mockRepo := new(UploadRepository)
	service := appUpload.NewUploadService(mockRepo)

	// Test data
	userID := int64(1)
	uuid := "not-found"

	// Mock expectations
	mockRepo.On("GetByUUID", mock.Anything, uuid, userID).
		Return(nil, appUpload.ErrFileNotFound)

	// Execute
	err := service.CanBeDeleted(context.Background(), uuid, userID)

	// Assert
	assert.Equal(t, appUpload.ErrFileNotFound, err)
	mockRepo.AssertExpectations(t)
}
