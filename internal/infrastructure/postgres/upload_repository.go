package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/aube/auth/internal/domain/entities"
	"github.com/aube/auth/internal/utils/logger"
	"github.com/rs/zerolog"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	queryUploadInsert           string = "INSERT INTO uploads (user_id, dataname, filename, description) VALUES ($1, $2, $3, $4) RETURNING id"
	queryUploadSelectByClientID string = "SELECT id, dataname, filename, description FROM uploads WHERE client_id = $1"
)

type UploadRepository struct {
	db  *pgxpool.Pool
	log zerolog.Logger
}

func NewUploadRepository(db *pgxpool.Pool) *UploadRepository {
	return &UploadRepository{
		db:  db,
		log: logger.Get().With().Str("postgres", "user_repository").Logger(),
	}
}

func (r *UploadRepository) Create(ctx context.Context, upload *entities.Upload) error {

	log.Print(upload)
	err := r.db.QueryRow(ctx, queryUploadInsert, 1, 2, 3, 4).Scan(&upload.ID)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UploadRepository) ListByUserID(ctx context.Context, username string) (*[]entities.Upload, error) {
	return nil, nil
}

func (r *UploadRepository) Delete(ctx context.Context, uuid string) error {
	return nil
}
