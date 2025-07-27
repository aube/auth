package handlers_upload

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aube/auth/internal/application/dto"
	"github.com/aube/auth/internal/domain/entities"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Custom response writer that implements CloseNotifier
type responseWriterCloseNotifier struct {
	*httptest.ResponseRecorder
	closeChan chan bool
}

func (w *responseWriterCloseNotifier) CloseNotify() <-chan bool {
	return w.closeChan
}

func (w *responseWriterCloseNotifier) Close() {
	w.closeChan <- true
}

// MockFileService реализует FileService интерфейс
type MockFileService struct {
	mock.Mock
}

func (m *MockFileService) Upload(ctx context.Context, size int64, data io.Reader) (*entities.File, error) {
	args := m.Called(ctx, size, data)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.File), args.Error(1)
}

func (m *MockFileService) Download(ctx context.Context, uuid string) (io.ReadCloser, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *MockFileService) Delete(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

// MockUploadService реализует UploadService интерфейс
type MockUploadService struct {
	mock.Mock
}

func (m *MockUploadService) RegisterUploadedFile(ctx context.Context, userID int64, file *entities.File, name, category, contentType, description string) (*entities.Upload, error) {
	args := m.Called(ctx, userID, file, name, category, contentType, description)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Upload), args.Error(1)
}

func (m *MockUploadService) GetByUUID(ctx context.Context, uuid string, userID int64) (*entities.Upload, error) {
	args := m.Called(ctx, uuid, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Upload), args.Error(1)
}

func (m *MockUploadService) GetByName(ctx context.Context, name string, userID int64) (*entities.Upload, error) {
	args := m.Called(ctx, name, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Upload), args.Error(1)
}

func (m *MockUploadService) ListByUserID(ctx context.Context, userID int64, offset, limit int) (*entities.Uploads, *dto.Pagination, error) {
	args := m.Called(ctx, userID, offset, limit)
	return args.Get(0).(*entities.Uploads), args.Get(1).(*dto.Pagination), args.Error(2)
}

func (m *MockUploadService) Delete(ctx context.Context, uuid string, userID int64) error {
	return m.Called(ctx, uuid, userID).Error(0)
}

func (m *MockUploadService) DeleteForce(ctx context.Context, uuid string, userID int64) error {
	return m.Called(ctx, uuid, userID).Error(0)
}
func TestUploadHandler_UploadFile_Success(t *testing.T) {
	// Setup
	mockFileService := new(MockFileService)
	mockUploadService := new(MockUploadService)
	handler := NewUploadHandler(mockFileService, mockUploadService)

	// Test data
	testContent := []byte("test content")
	expectedFile := &entities.File{
		Name: "test-uuid",
		Size: int64(len(testContent)),
	}
	expectedUpload := &entities.Upload{
		UUID:        "test-uuid",
		Name:        "test.txt",
		ContentType: "text/plain",
		Size:        int64(len(testContent)),
	}

	// Mock expectations
	mockUploadService.On("GetByName", mock.Anything, "test.txt", int64(1)).
		Return(nil, errors.New("not found"))

	// Важное изменение: используем mock.MatchedBy для проверки reader
	mockFileService.On("Upload",
		mock.Anything,           // context
		int64(len(testContent)), // size
		mock.MatchedBy(func(r io.Reader) bool {
			data, err := io.ReadAll(r)
			return err == nil && string(data) == string(testContent)
		}),
	).Return(expectedFile, nil)

	mockUploadService.On("RegisterUploadedFile",
		mock.Anything,
		int64(1),
		expectedFile,
		"test.txt",
		"docs",
		mock.MatchedBy(func(contentType string) bool {
			return contentType == "text/plain" || contentType == "application/octet-stream"
		}),
		"test file",
	).Return(expectedUpload, nil)

	// Create test request
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", "test.txt")
	require.NoError(t, err)
	_, err = part.Write(testContent)
	require.NoError(t, err)

	err = writer.WriteField("description", "test file")
	require.NoError(t, err)
	err = writer.WriteField("category", "docs")
	require.NoError(t, err)
	writer.Close()

	// Create gin context with custom writer
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", int64(1))
	c.Request = httptest.NewRequest("POST", "/upload", body)
	c.Request.Header.Set("Content-Type", writer.FormDataContentType())

	// Execute
	handler.UploadFile(c)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), `"uuid":"test-uuid"`)

	mockFileService.AssertExpectations(t)
	mockUploadService.AssertExpectations(t)
}

func TestUploadHandler_DownloadFile_Success(t *testing.T) {
	// Setup
	mockFileService := new(MockFileService)
	mockUploadService := new(MockUploadService)
	handler := NewUploadHandler(mockFileService, mockUploadService)

	// Mock data
	uuid := "test-uuid"
	expectedUpload := &entities.Upload{
		UUID:        uuid,
		Name:        "test.txt",
		ContentType: "text/plain",
		Size:        123,
	}
	fileContent := io.NopCloser(bytes.NewReader([]byte("test content")))

	// Mock expectations
	mockUploadService.On("GetByUUID", mock.Anything, uuid, int64(1)).Return(expectedUpload, nil)
	mockFileService.On("Download", mock.Anything, uuid).Return(fileContent, nil)

	// Create a custom response writer that implements CloseNotifier
	w := &responseWriterCloseNotifier{
		ResponseRecorder: httptest.NewRecorder(),
	}

	// Create gin context with our custom writer
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", int64(1))
	c.Request = httptest.NewRequest("GET", "/download?uuid="+uuid, nil)

	// Execute
	handler.DownloadFile(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/plain", w.Header().Get("Content-Type"))
	assert.Equal(t, "attachment; filename=test.txt", w.Header().Get("Content-Disposition"))
	assert.Equal(t, "123", w.Header().Get("Content-Length"))
	mockFileService.AssertExpectations(t)
	mockUploadService.AssertExpectations(t)
}

func TestUploadHandler_ListFiles_Success(t *testing.T) {
	// Setup
	mockFileService := new(MockFileService)
	mockUploadService := new(MockUploadService)
	handler := NewUploadHandler(mockFileService, mockUploadService)

	// Mock expectations
	uploads := &entities.Uploads{
		entities.Upload{UUID: "1", Name: "file1.txt"},
		entities.Upload{UUID: "2", Name: "file2.txt"},
	}
	pagination := &dto.Pagination{Total: 2, Size: 10, Page: 1}

	mockUploadService.On("ListByUserID", mock.Anything, int64(1), 0, 10).Return(uploads, pagination, nil)

	// Create test context
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("userID", int64(1))
	c.Set("offset", 0)
	c.Set("limit", 10)
	c.Request = httptest.NewRequest("GET", "/uploads", nil)

	// Test
	handler.ListFiles(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"rows":`)
	assert.Contains(t, w.Body.String(), `"pagination":`)
	mockUploadService.AssertExpectations(t)
}

func TestUploadHandler_DeleteFile_Success(t *testing.T) {
	// Setup
	mockFileService := new(MockFileService)
	mockUploadService := new(MockUploadService)
	handler := NewUploadHandler(mockFileService, mockUploadService)

	// Mock expectations
	uuid := "aaaaaaaa-aaaa-bbbb-cccc-aaaabbbbcccc"
	mockUploadService.On("Delete", mock.Anything, uuid, int64(1)).Return(nil)
	mockFileService.On("Delete", mock.Anything, uuid).Return(nil)

	// Create proper router for the test
	router := gin.Default()
	router.DELETE("/file", func(c *gin.Context) {
		c.Set("userID", int64(1))
		handler.DeleteFile(c)
	})

	// Create test request
	w := httptest.NewRecorder()
	req := httptest.NewRequest("DELETE", "/file?uuid="+uuid, nil)

	// Execute request
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String()) // Проверяем что тело ответа пустое
	mockFileService.AssertExpectations(t)
	mockUploadService.AssertExpectations(t)
}
