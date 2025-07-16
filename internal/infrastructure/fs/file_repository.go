package fs

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"sync"

	appFile "github.com/aube/auth/internal/application/file"
	"github.com/aube/auth/internal/domain/entities"
	"github.com/aube/auth/internal/utils/logger"
	"github.com/rs/zerolog"
)

// FileSystemRepository реализация FileRepository для хранения файлов в файловой системе
type FileSystemRepository struct {
	storagePath string
	mu          sync.RWMutex
	log         zerolog.Logger
}

// NewFileSystemRepository создает новый экземпляр FileSystemRepository
func NewFileSystemRepository(storagePath string) (*FileSystemRepository, error) {
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		return nil, err
	}
	return &FileSystemRepository{
		storagePath: storagePath,
		log:         logger.Get().With().Str("fs", "file_repository").Logger(),
	}, nil
}

func (r *FileSystemRepository) Save(ctx context.Context, file *entities.File, data io.Reader) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	filePath := filepath.Join(r.storagePath, file.ID)
	dst, err := os.Create(filePath)
	if err != nil {
		r.log.Debug().Err(err).Msg("Save")
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, data); err != nil {
		r.log.Debug().Err(err).Msg("Save")
		return err
	}

	return nil
}

func (r *FileSystemRepository) FindByID(ctx context.Context, id string) (*entities.File, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	filePath := filepath.Join(r.storagePath, id)
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		r.log.Debug().Err(err).Msg("FindByID")
		if os.IsNotExist(err) {
			return nil, appFile.ErrFileNotFound
		}
		return nil, err
	}

	// В реальном приложении нужно хранить метаданные в БД
	// Здесь упрощенная реализация
	return &entities.File{
		ID:         id,
		Name:       id, // В реальном приложении имя должно храниться отдельно
		Size:       fileInfo.Size(),
		Path:       filePath,
		UploadedAt: fileInfo.ModTime(),
	}, nil
}

func (r *FileSystemRepository) FindAll(ctx context.Context) ([]*entities.File, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	files, err := os.ReadDir(r.storagePath)
	if err != nil {
		r.log.Debug().Err(err).Msg("FindAll")
		return nil, err
	}

	var result []*entities.File
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileInfo, err := file.Info()
		if err != nil {
			continue
		}

		result = append(result, &entities.File{
			ID:         file.Name(),
			Name:       file.Name(),
			Size:       fileInfo.Size(),
			Path:       filepath.Join(r.storagePath, file.Name()),
			UploadedAt: fileInfo.ModTime(),
		})
	}

	return result, nil
}

func (r *FileSystemRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	filePath := filepath.Join(r.storagePath, id)
	if err := os.Remove(filePath); err != nil {
		r.log.Debug().Err(err).Msg("Delete")
		if os.IsNotExist(err) {
			return appFile.ErrFileNotFound
		}
		return err
	}
	return nil
}

func (r *FileSystemRepository) GetFileContent(ctx context.Context, file *entities.File) (io.ReadCloser, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return os.Open(file.Path)
}

var _ appFile.FileRepository = (*FileSystemRepository)(nil)
