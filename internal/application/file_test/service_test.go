package file_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/aube/auth/internal/application/file"
	"github.com/aube/auth/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFileService_Upload_Success(t *testing.T) {
	// Setup
	mockRepo := new(FileRepository)
	service := file.NewFileService(mockRepo)

	// Test data
	testContent := []byte("test file content")

	// Mock expectations
	mockRepo.On("Save", mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			fileArg := args.Get(1).(*entities.File)
			assert.NotEmpty(t, fileArg.Name)
			assert.Equal(t, int64(len(testContent)), fileArg.Size)

			// Проверяем содержимое файла
			data, err := io.ReadAll(args.Get(2).(io.Reader))
			assert.NoError(t, err)
			assert.Equal(t, testContent, data)
		}).
		Return(nil)

	// Execute
	file, err := service.Upload(context.Background(), int64(len(testContent)), bytes.NewReader(testContent))

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, file)
	assert.Equal(t, int64(len(testContent)), file.Size)
	mockRepo.AssertExpectations(t)
}

func TestFileService_Upload_RepositoryError(t *testing.T) {
	// Setup
	mockRepo := new(FileRepository)
	service := file.NewFileService(mockRepo)

	// Mock expectations
	expectedError := errors.New("repository error")
	mockRepo.On("Save", mock.Anything, mock.Anything, mock.Anything).
		Return(expectedError)

	// Execute
	file, err := service.Upload(context.Background(), 123, bytes.NewReader([]byte("test")))

	// Assert
	assert.Nil(t, file)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

func TestFileService_Download_Success(t *testing.T) {
	// Setup
	mockRepo := new(FileRepository)
	service := file.NewFileService(mockRepo)

	// Test data
	expectedContent := []byte("test content")
	uuid := "test-uuid"

	// Mock expectations
	mockRepo.On("GetFileContent", mock.Anything, uuid).
		Return(io.NopCloser(bytes.NewReader(expectedContent)), nil)

	// Execute
	reader, err := service.Download(context.Background(), uuid)

	// Assert
	require.NoError(t, err)
	defer reader.Close()

	data, err := io.ReadAll(reader)
	require.NoError(t, err)
	assert.Equal(t, expectedContent, data)
	mockRepo.AssertExpectations(t)
}

func TestFileService_Download_NotFound(t *testing.T) {
	// Setup
	mockRepo := new(FileRepository)
	service := file.NewFileService(mockRepo)

	// Mock expectations
	mockRepo.On("GetFileContent", mock.Anything, "not-found").
		Return(nil, file.ErrFileNotFound)

	// Execute
	reader, err := service.Download(context.Background(), "not-found")

	// Assert
	assert.Nil(t, reader)
	assert.Equal(t, file.ErrFileNotFound, err)
	mockRepo.AssertExpectations(t)
}

func TestFileService_Delete_Success(t *testing.T) {
	// Setup
	mockRepo := new(FileRepository)
	service := file.NewFileService(mockRepo)

	// Mock expectations
	mockRepo.On("Delete", mock.Anything, "test-uuid").Return(nil)

	// Execute
	err := service.Delete(context.Background(), "test-uuid")

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestFileService_Delete_Error(t *testing.T) {
	// Setup
	mockRepo := new(FileRepository)
	service := file.NewFileService(mockRepo)

	// Mock expectations
	expectedError := errors.New("delete error")
	mockRepo.On("Delete", mock.Anything, "test-uuid").Return(expectedError)

	// Execute
	err := service.Delete(context.Background(), "test-uuid")

	// Assert
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}
