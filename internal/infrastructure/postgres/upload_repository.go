package postgres

import (
	"context"
	"fmt"
	"strconv"

	"github.com/aube/auth/internal/domain/entities"
	"github.com/aube/auth/internal/utils/logger"
	"github.com/rs/zerolog"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	queryUploadInsert           string = "INSERT INTO uploads (user_id, name, uuid, size, type, description) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id"
	queryUploadSelectByClientID string = "SELECT id, dataname, filename, description FROM uploads WHERE client_id = $1"
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

func (r *UploadRepository) Create(ctx context.Context, userID string, uploadedFile *entities.File) error {

	usid, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		r.log.Debug().Err(err).Msg("Create1")
		return fmt.Errorf("failed to create upload: %w", err)
	}

	r.log.Debug().Msg("Upload")
	r.log.Debug().Msg(uploadedFile.Name)
	r.log.Debug().Msg(uploadedFile.UUID)

	err = r.db.QueryRow(ctx,
		queryUploadInsert,
		usid,
		uploadedFile.Name,
		uploadedFile.UUID,
		uploadedFile.Size,
		uploadedFile.ContentType,
		uploadedFile.Description,
	).Scan(&uploadedFile.ID)
	if err != nil {
		r.log.Debug().Err(err).Msg("Create2")
		return fmt.Errorf("failed to create upload: %w", err)
	}

	return nil

}

func (r *UploadRepository) ListByUserID(ctx context.Context, username string) (*[]entities.Upload, error) {
	return nil, nil
}

func (r *UploadRepository) Delete(ctx context.Context, uuid string) error {
	return nil
}
