// Package upload provides data persistence operations for file upload metadata.
package upload

import (
	"context"
	"errors"

	"github.com/aube/auth/internal/application/dto"
	"github.com/aube/auth/internal/domain/entities"
)

// ErrFileNotFound is returned when requested upload metadata cannot be found.
var ErrFileNotFound = errors.New("file not found")

// UploadRepository defines the interface for upload metadata persistence operations.
// Implementations should handle database operations for upload records.
//
// Methods:
//
//   - Create: Stores new upload metadata
//     ctx: Context for cancellation/timeout
//     userID: Owner of the upload
//     upload: Upload metadata entity
//     Returns: error on failure
//
//   - ListByUserID: Retrieves paginated uploads for a user
//     ctx: Context for cancellation/timeout
//     userID: Owner filter
//     offset: Pagination offset
//     limit: Maximum results per page
//     Returns: (*entities.Uploads, *dto.Pagination, error)
//
//   - GetByUUID: Retrieves upload by unique identifier with owner check
//     ctx: Context for cancellation/timeout
//     uuid: Upload identifier
//     userID: Owner verification
//     Returns: (*entities.Upload, error)
//
//   - GetByName: Retrieves upload by filename with owner check
//     ctx: Context for cancellation/timeout
//     name: Filename
//     userID: Owner verification
//     Returns: (*entities.Upload, error)
//
//   - Delete: Soft-deletes upload record (standard deletion)
//     ctx: Context for cancellation/timeout
//     uuid: Upload identifier
//     userID: Owner verification
//     Returns: error on failure
//
//   - DeleteForce: Hard-deletes upload record (bypasses normal checks)
//     ctx: Context for cancellation/timeout
//     uuid: Upload identifier
//     userID: Owner verification
//     Returns: error on failure
type UploadRepository interface {
	Create(ctx context.Context, userID int64, upload *entities.Upload) error
	ListByUserID(ctx context.Context, userID int64, offset, limit int, params map[string]any) (*entities.Uploads, *dto.Pagination, error)
	GetByUUID(ctx context.Context, uuid string, userID int64) (*entities.Upload, error)
	GetByName(ctx context.Context, name string, userID int64) (*entities.Upload, error)
	Delete(ctx context.Context, uuid string, userID int64) error
	DeleteForce(ctx context.Context, uuid string, userID int64) error
}
