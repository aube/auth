package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	appUpload "github.com/aube/auth/internal/application/upload"
	"github.com/aube/auth/internal/domain/entities"
	"github.com/aube/auth/internal/utils/logger"
	"github.com/rs/zerolog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	queryUploadInsert         string = "INSERT INTO uploads (user_id, uuid, size, name, content_type, description) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	queryUploadSelectByUserID string = "SELECT id, user_id, uuid, size, name, contentType, description, createdAt} FROM uploads WHERE user_id = $1"
	queryGetByUUID            string = "SELECT id, user_id, size, name, content_type, description, created_at) FROM uploads WHERE uuid = $1 and user_id=$2"
)

type UploadRepository struct {
	db  *pgxpool.Pool
	log zerolog.Logger
}

func NewUploadRepository(db *pgxpool.Pool) *UploadRepository {
	return &UploadRepository{
		db:  db,
		log: logger.Get().With().Str("postgres", "upload_repository").Logger(),
	}
}

func (r *UploadRepository) Create(
	ctx context.Context,
	userID int64,
	upload *entities.Upload,
) error {
	r.log.Debug().Msg("Upload")
	r.log.Debug().Msg(upload.Name)
	r.log.Debug().Msg(upload.UUID)

	var id int64
	err := r.db.QueryRow(ctx,
		queryUploadInsert,
		userID,
		upload.UUID,
		upload.Size,
		upload.Name,
		upload.ContentType,
		upload.Description,
	).Scan(&id)

	if err != nil {
		r.log.Debug().Err(err).Msg("Create2")
		return fmt.Errorf("failed to create upload: %w", err)
	}

	return nil
}

// List returns all URL mappings for the current user from the database.
// Returns an unauthorized error if no user ID is present in context.
func (r *UploadRepository) ListByUserID(ctx context.Context, userID int64) (*entities.Uploads, error) {

	rows, err := r.db.Query(ctx, queryUploadSelectByUserID, userID)
	if err != nil {
		r.log.Debug().Err(err).Msg("ListByUserID1")
		return nil, err
	}

	if err = rows.Err(); err != nil {
		r.log.Debug().Err(err).Msg("ListByUserID2")
		return nil, err
	}

	var uploads entities.Uploads

	var (
		id          int64
		uuid        string
		user_id     int64
		size        int64
		name        string
		contentType string
		description string
		createdAt   time.Time
	)

	pgx.ForEachRow(rows, []any{&id, &uuid, &user_id, &size, &name, &contentType, &description, &createdAt}, func() error {
		file := entities.NewFile(uuid, "", size)
		upload := entities.NewUpload(
			file,
			id,
			user_id,
			name,
			contentType,
			description,
			createdAt,
		)
		uploads = append(uploads, *upload)
		return nil
	})

	return &uploads, nil
}

// List returns all URL mappings for the current user from the database.
// Returns an unauthorized error if no user ID is present in context.
func (r *UploadRepository) GetByUUID(ctx context.Context, uuid string, userID int64) (*entities.Upload, error) {
	var (
		id          int64
		user_id     int64
		name        string
		size        int64
		contentType string
		description string
		createdAt   time.Time
	)

	err := r.db.QueryRow(ctx, queryGetByUUID, uuid, userID).Scan(&id, &user_id, &name, &size, &contentType, &description, &createdAt)

	if err != nil {
		r.log.Debug().Err(err).Msg("GetByUUID")
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appUpload.ErrFileNotFound
		}
		return nil, fmt.Errorf("failed to find upload: %w", err)
	}

	file := entities.NewFile(uuid, "", size)

	return entities.NewUpload(
		file,
		id,
		user_id,
		name,
		contentType,
		description,
		createdAt,
	), nil
}
