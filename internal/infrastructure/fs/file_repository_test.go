package fs

import (
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/aube/auth/internal/application/file"
	"github.com/aube/auth/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestFS(t *testing.T) (string, *FileSystemRepository, func()) {
	tempDir, err := os.MkdirTemp("", "filestorage_test")
	require.NoError(t, err)

	repo, err := NewFileSystemRepository(tempDir)
	require.NoError(t, err)

	cleanup := func() {
		os.RemoveAll(tempDir)
	}

	return tempDir, repo, cleanup
}

func TestNewFileSystemRepository(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tempDir := filepath.Join(os.TempDir(), "new_repo_test")
		defer os.RemoveAll(tempDir)

		repo, err := NewFileSystemRepository(tempDir)
		require.NoError(t, err)
		assert.NotNil(t, repo)
		assert.DirExists(t, tempDir)
	})

	t.Run("invalid path", func(t *testing.T) {
		// Попытка создать в несуществующем корневом каталоге
		repo, err := NewFileSystemRepository("/nonexistent/path/to/repo")
		assert.Error(t, err)
		assert.Nil(t, repo)
	})
}

func TestFileSystemRepository_Save(t *testing.T) {
	_, repo, cleanup := setupTestFS(t)
	defer cleanup()

	testContent := []byte("test file content")
	testFile := &entities.File{
		Name: "testfile.txt",
		Size: int64(len(testContent)),
	}

	t.Run("success", func(t *testing.T) {
		err := repo.Save(context.Background(), testFile, bytes.NewReader(testContent))
		require.NoError(t, err)

		// Проверяем что файл действительно создался
		filePath := filepath.Join(repo.storagePath, testFile.Name)
		fileInfo, err := os.Stat(filePath)
		require.NoError(t, err)
		assert.Equal(t, testFile.Name, fileInfo.Name())
		assert.Equal(t, testFile.Size, fileInfo.Size())
	})

	t.Run("empty reader", func(t *testing.T) {
		err := repo.Save(context.Background(), testFile, bytes.NewReader(nil))
		require.NoError(t, err)
	})

	t.Run("concurrent save", func(t *testing.T) {
		concurrentFile := &entities.File{Name: "concurrent.txt", Size: 10}
		content := []byte("content")

		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err := repo.Save(context.Background(), concurrentFile, bytes.NewReader(content))
				assert.NoError(t, err)
			}()
		}
		wg.Wait()

		// Должен быть сохранен только один файл
		filePath := filepath.Join(repo.storagePath, concurrentFile.Name)
		_, err := os.Stat(filePath)
		assert.NoError(t, err)
	})
}

func TestFileSystemRepository_FindAll(t *testing.T) {
	tempDir, repo, cleanup := setupTestFS(t)
	defer cleanup()

	// Создаем тестовые файлы
	files := []struct {
		name    string
		content string
	}{
		{"file1.txt", "content1"},
		{"file2.txt", "content2"},
	}

	for _, f := range files {
		path := filepath.Join(tempDir, f.name)
		err := os.WriteFile(path, []byte(f.content), 0644)
		require.NoError(t, err)
	}

	t.Run("success", func(t *testing.T) {
		result, err := repo.FindAll(context.Background())
		require.NoError(t, err)
		require.Len(t, *result, 2)

		// Проверяем что файлы в правильном порядке и с правильными данными
		assert.Equal(t, "file1.txt", (*result)[0].Name)
		assert.Equal(t, int64(len("content1")), (*result)[0].Size)
		assert.Equal(t, "file2.txt", (*result)[1].Name)
		assert.Equal(t, int64(len("content2")), (*result)[1].Size)
	})

	t.Run("empty storage", func(t *testing.T) {
		_, emptyRepo, emptyCleanup := setupTestFS(t)
		defer emptyCleanup()

		result, err := emptyRepo.FindAll(context.Background())
		require.NoError(t, err)
		assert.Empty(t, *result)
	})
}

func TestFileSystemRepository_Delete(t *testing.T) {
	tempDir, repo, cleanup := setupTestFS(t)
	defer cleanup()

	// Создаем тестовый файл
	testFile := "test_delete.txt"
	filePath := filepath.Join(tempDir, testFile)
	err := os.WriteFile(filePath, []byte("content"), 0644)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		err := repo.Delete(context.Background(), testFile)
		assert.NoError(t, err)
		assert.NoFileExists(t, filePath)
	})

	t.Run("not found", func(t *testing.T) {
		err := repo.Delete(context.Background(), "nonexistent.txt")
		assert.Equal(t, file.ErrFileNotFound, err)
	})
}

func TestFileSystemRepository_GetFileContent(t *testing.T) {
	tempDir, repo, cleanup := setupTestFS(t)
	defer cleanup()

	// Создаем тестовый файл
	testFile := "test_content.txt"
	testContent := []byte("test file content")
	filePath := filepath.Join(tempDir, testFile)
	err := os.WriteFile(filePath, testContent, 0644)
	require.NoError(t, err)

	t.Run("success", func(t *testing.T) {
		reader, err := repo.GetFileContent(context.Background(), testFile)
		require.NoError(t, err)
		defer reader.Close()

		content, err := io.ReadAll(reader)
		require.NoError(t, err)
		assert.Equal(t, testContent, content)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.GetFileContent(context.Background(), "nonexistent.txt")
		assert.Equal(t, file.ErrFileNotFound, err)
	})

	t.Run("concurrent read", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				reader, err := repo.GetFileContent(context.Background(), testFile)
				require.NoError(t, err)
				defer reader.Close()

				content, err := io.ReadAll(reader)
				require.NoError(t, err)
				assert.Equal(t, testContent, content)
			}()
		}
		wg.Wait()
	})
}

func TestFileSystemRepository_InterfaceImplementation(t *testing.T) {
	_, repo, cleanup := setupTestFS(t)
	defer cleanup()

	// Проверяем что репозиторий реализует интерфейс FileRepository
	var _ file.FileRepository = repo
}
