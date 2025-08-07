package postgres

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/aube/auth/internal/application/dto"
	appUpload "github.com/aube/auth/internal/application/upload"
	"github.com/aube/auth/internal/domain/entities"
	"github.com/aube/auth/internal/utils/logger"
	"github.com/aube/auth/internal/utils/sql"
	"github.com/rs/zerolog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	queryImageInsert              string = "INSERT INTO images (user_id, uuid, size, name, category, content_type, description) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id"
	queryImageSelectByUserID      string = "SELECT id, user_id, uuid, size, name, category, content_type, description, created_at FROM images %WHERE% OFFSET $1 LIMIT $2"
	queryImageSelectByUserIDTotal string = "SELECT count(*) total FROM images %WHERE%"
	queryImageGetByUUID           string = "SELECT id, user_id, size, name, category, content_type, description, created_at FROM images WHERE uuid = $1 and user_id=$2 and deleted=false"
	queryImageGetByName           string = "SELECT id, user_id, uuid, size, category, content_type, description, created_at FROM images WHERE name = $1 and user_id=$2 and deleted=false"
	queryImageDelete              string = "UPDATE images SET deleted=true WHERE uuid = $1 and user_id=$2"
	queryImageDeleteForce         string = "DELETE images WHERE uuid = $1 and user_id=$2"
)

type ImageRepository struct {
	db  *pgxpool.Pool
	log zerolog.Logger
}

func NewImageRepository(db *pgxpool.Pool) *ImageRepository {
	return &ImageRepository{
		db:  db,
		log: logger.Get().With().Str("postgres", "image_repository").Logger(),
	}
}

func (r *ImageRepository) Create(
	ctx context.Context,
	userID int64,
	image *entities.Image,
) error {
	r.log.Debug().Msg("Image image")
	r.log.Debug().Msg(image.Name)
	r.log.Debug().Msg(image.UUID)

	var id int64
	err := r.db.QueryRow(ctx,
		queryImageInsert,
		userID,
		image.UUID,
		image.Size,
		image.Name,
		image.Category,
		image.ContentType,
		image.Description,
	).Scan(&id)

	if err != nil {
		r.log.Debug().Err(err).Msg("Create2")
		return fmt.Errorf("failed to create image: %w", err)
	}

	return nil
}

// List returns all URL mappings for the current user from the database.
// Returns an unauthorized error if no user ID is present in context.
func (r *ImageRepository) ListByUserID(ctx context.Context, userID int64, offset, limit int, params map[string]any) (*entities.Images, *dto.Pagination, error) {

	whereClause, whereParams := sql.BuildWhere(params, "AND", 3)
	allParams := []any{offset, limit}
	allParams = append(allParams, whereParams...)

	query := strings.Replace(queryImageSelectByUserID, "%WHERE%", "WHERE "+whereClause, 1)

	rows, err := r.db.Query(ctx, query, allParams...)
	if err != nil {
		r.log.Debug().Err(err).Msg("ListByUserID1")
		return nil, nil, err
	}
	defer rows.Close()

	var images entities.Images

	for rows.Next() {
		var (
			id          int64
			userId      int64
			uuid        string
			size        int64
			name        string
			category    string
			contentType string
			description string
			createdAt   time.Time
		)

		err := rows.Scan(
			&id,
			&userId,
			&uuid,
			&size,
			&name,
			&category,
			&contentType,
			&description,
			&createdAt,
		)
		if err != nil {
			r.log.Error().Err(err).Msg("Failed to scan image row")
			return nil, nil, fmt.Errorf("failed to scan image row: %w", err)
		}

		file := entities.NewFile(uuid, "", size)
		image := entities.NewImage(
			file,
			id,
			userId,
			name,
			category,
			contentType,
			description,
			createdAt,
		)
		images = append(images, *image)
	}

	if err = rows.Err(); err != nil {
		r.log.Debug().Err(err).Msg("ListByUserID2")
		return nil, nil, fmt.Errorf("error after iterating image rows: %w", err)
	}

	// Totals
	whereClause, whereParams = sql.BuildWhere(params, "AND", 1)
	query = strings.Replace(queryImageSelectByUserIDTotal, "%WHERE%", "WHERE "+whereClause, 1)

	var total int
	err = r.db.QueryRow(ctx, query, whereParams...).Scan(&total)

	if err != nil {
		r.log.Debug().Err(err).Msg("GetTotals")
		return nil, nil, fmt.Errorf("failed to get totals: %w", err)
	}

	page := float64(offset) / float64(limit)
	pagination := dto.Pagination{
		Total: total,
		Page:  int(math.Round(page)) + 1,
		Size:  limit,
	}

	return &images, &pagination, nil
}

func (r *ImageRepository) GetByUUID(ctx context.Context, uuid string, userID int64) (*entities.Image, error) {
	var (
		id          int64
		user_id     int64
		name        string
		category    string
		size        int64
		contentType string
		description string
		createdAt   time.Time
	)

	err := r.db.QueryRow(ctx, queryImageGetByUUID, uuid, userID).Scan(&id, &user_id, &size, &name, &category, &contentType, &description, &createdAt)

	if err != nil {
		r.log.Debug().Err(err).Msg("GetByUUID")
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appUpload.ErrFileNotFound
		}
		return nil, fmt.Errorf("failed to find image: %w", err)
	}

	file := entities.NewFile(uuid, "", size)

	return entities.NewImage(
		file,
		id,
		user_id,
		name,
		category,
		contentType,
		description,
		createdAt,
	), nil
}

func (r *ImageRepository) GetByName(ctx context.Context, name string, userID int64) (*entities.Image, error) {
	var (
		id          int64
		user_id     int64
		uuid        string
		size        int64
		category    string
		contentType string
		description string
		createdAt   time.Time
	)

	err := r.db.QueryRow(ctx, queryImageGetByName, name, userID).Scan(&id, &user_id, &uuid, &size, &category, &contentType, &description, &createdAt)

	if err != nil {
		r.log.Debug().Err(err).Msg("GetByName")
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appUpload.ErrFileNotFound
		}
		return nil, fmt.Errorf("failed to find image: %w", err)
	}

	file := entities.NewFile(uuid, "", size)

	return entities.NewImage(
		file,
		id,
		user_id,
		name,
		category,
		contentType,
		description,
		createdAt,
	), nil
}

func (r *ImageRepository) Delete(ctx context.Context, uuid string, userID int64) error {
	_, err := r.db.Query(ctx, queryImageDelete, uuid, userID)

	if err != nil {
		r.log.Debug().Err(err).Msg("Delete")
		if errors.Is(err, pgx.ErrNoRows) {
			return appUpload.ErrFileNotFound
		}
		return fmt.Errorf("failed to delete image: %w", err)
	}

	return nil
}

func (r *ImageRepository) DeleteForce(ctx context.Context, uuid string, userID int64) error {
	_, err := r.db.Query(ctx, queryImageDeleteForce, uuid, userID)

	if err != nil {
		r.log.Debug().Err(err).Msg("Delete")
		if errors.Is(err, pgx.ErrNoRows) {
			return appUpload.ErrFileNotFound
		}
		return fmt.Errorf("failed to delete image: %w", err)
	}

	return nil
}
