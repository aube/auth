package postgres

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/aube/auth/internal/application/dto"
	appPage "github.com/aube/auth/internal/application/page"
	"github.com/aube/auth/internal/domain/entities"
	"github.com/aube/auth/internal/utils/logger"
	"github.com/aube/auth/internal/utils/sql"
	"github.com/rs/zerolog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	pageFields            string = "name, meta, title, category, template, h1, content, content_short, created_at, updated_at"
	queryPageInsert       string = "INSERT INTO pages (" + pageFields + ") VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id"
	queryPageUpdate       string = "UPDATE pages SET name=$1, meta=$2, title=$3, category=$4, template=$5, h1=$6, content=$7, content_short=$8 WHERE id=$9"
	queryPageSelectByName string = "SELECT id, " + pageFields + " FROM pages WHERE name = $1 and deleted = false"
	queryPageSelectByID   string = "SELECT id, " + pageFields + " FROM pages WHERE id = $1 and deleted = false"
	queryPageGetIDByName  string = "SELECT id FROM pages WHERE name = $1"
	queryPageDelete       string = "UPDATE pages SET deleted=true WHERE id = $1"
	queryPageDeleteForce  string = "DELETE FROM pages WHERE id = $1"
	queryPagesSelect      string = "SELECT id, " + pageFields + " FROM pages %WHERE% OFFSET $1 LIMIT $2"
	queryPagesSelectTotal string = "SELECT count(*) total FROM pages %WHERE%"
)

type PageRepository struct {
	db  *pgxpool.Pool
	log zerolog.Logger
}

func NewPageRepository(db *pgxpool.Pool) *PageRepository {
	return &PageRepository{
		db:  db,
		log: logger.Get().With().Str("postgres", "page_repository").Logger(),
	}
}

func (r *PageRepository) Create(ctx context.Context, page *entities.Page) error {
	err := r.db.QueryRow(
		ctx,
		queryPageInsert,
		page.Name,
		page.Meta,
		page.Title,
		page.Category,
		page.Template,
		page.H1,
		page.Content,
		page.ContentShort,
	).Scan(&page.ID)

	if err != nil {
		r.log.Debug().Err(err).Msg(page.Name)
		r.log.Debug().Err(err).Msg("Create")
		return fmt.Errorf("failed to create page: %w", err)
	}

	return nil
}

func (r *PageRepository) Update(ctx context.Context, page *entities.Page) error {
	_, err := r.db.Exec(
		ctx,
		queryPageUpdate,
		page.Name,
		page.Meta,
		page.Title,
		page.Category,
		page.Template,
		page.H1,
		page.Content,
		page.ContentShort,
		page.ID,
	)

	if err != nil {
		r.log.Debug().Err(err).Msg(page.Name)
		r.log.Debug().Err(err).Msg("Update")
		return fmt.Errorf("failed to update page: %w", err)
	}

	return nil
}

func (r *PageRepository) FindByName(ctx context.Context, param string) (*entities.PageWithTime, error) {
	var (
		id           int64
		name         string
		meta         string
		title        string
		category     string
		template     string
		h1           string
		content      string
		contentShort string
		createdAt    time.Time
		updatedAt    time.Time
	)
	err := r.db.QueryRow(
		ctx,
		queryPageSelectByName,
		param,
	).Scan(
		&id,
		&name,
		&meta,
		&title,
		&category,
		&template,
		&h1,
		&content,
		&contentShort,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		r.log.Debug().Err(err).Msg("FindByName")
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appPage.ErrPageNotFound
		}
		return nil, fmt.Errorf("failed to find page: %w", err)
	}

	return entities.NewPageWithTime(
		id,
		name,
		meta,
		title,
		category,
		template,
		h1,
		content,
		contentShort,
		createdAt,
		updatedAt,
	)
}

func (r *PageRepository) FindByID(ctx context.Context, param int64) (*entities.PageWithTime, error) {
	var (
		id           int64
		name         string
		meta         string
		title        string
		category     string
		template     string
		h1           string
		content      string
		contentShort string
		createdAt    time.Time
		updatedAt    time.Time
	)
	err := r.db.QueryRow(
		ctx,
		queryPageSelectByID,
		param,
	).Scan(
		&id,
		&name,
		&meta,
		&title,
		&category,
		&template,
		&h1,
		&content,
		&contentShort,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		r.log.Debug().Err(err).Msg("FindByName")
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, appPage.ErrPageNotFound
		}
		return nil, fmt.Errorf("failed to find page: %w", err)
	}

	return entities.NewPageWithTime(
		id,
		name,
		meta,
		title,
		category,
		template,
		h1,
		content,
		contentShort,
		createdAt,
		updatedAt,
	)
}

func (r *PageRepository) GetIDByName(ctx context.Context, name string) (int64, error) {
	var id int64
	err := r.db.QueryRow(ctx, queryPageGetIDByName, name).Scan(&id)

	if err != nil {
		if err.Error() == "no rows in result set" {
			return 0, nil
		}

		r.log.Debug().Err(err).Msg("Exists")
		return 0, fmt.Errorf("failed to check page existence: %w", err)
	}

	return id, nil
}

func (r *PageRepository) Delete(ctx context.Context, pageID int64) error {

	_, err := r.db.Query(ctx, queryPageDelete, pageID)

	if err != nil {
		r.log.Debug().Err(err).Msg("Delete")
		return fmt.Errorf("failed to delete page: %w", err)
	}

	return nil
}

func (r *PageRepository) DeleteForce(ctx context.Context, pageID int64) error {

	_, err := r.db.Query(ctx, queryPageDeleteForce, pageID)

	if err != nil {
		r.log.Debug().Err(err).Msg("Delete")
		return fmt.Errorf("failed to delete page: %w", err)
	}

	return nil
}

// List returns all Pages from the database.
func (r *PageRepository) ListPages(ctx context.Context, offset, limit int, params map[string]any) (*entities.PagesWithTimes, *dto.Pagination, error) {

	whereClause, whereParams := sql.BuildWhere(params, "AND", 3)
	allParams := []any{offset, limit}
	allParams = append(allParams, whereParams...)

	query := strings.Replace(queryPagesSelect, "%WHERE%", "WHERE "+whereClause, 1)

	rows, err := r.db.Query(ctx, query, allParams...)
	if err != nil {
		r.log.Debug().Err(err).Msg("ListPages1")
		return nil, nil, err
	}
	defer rows.Close()

	var pages entities.PagesWithTimes

	for rows.Next() {
		var (
			id           int64
			name         string
			meta         string
			title        string
			category     string
			template     string
			h1           string
			content      string
			contentShort string
			createdAt    time.Time
			updatedAt    time.Time
		)

		err := rows.Scan(
			&id,
			&name,
			&meta,
			&title,
			&category,
			&template,
			&h1,
			&content,
			&contentShort,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			r.log.Error().Err(err).Msg("Failed to scan upload row")
			return nil, nil, fmt.Errorf("failed to scan upload row: %w", err)
		}

		page, _ := entities.NewPageWithTime(
			id,
			name,
			meta,
			title,
			category,
			template,
			h1,
			content,
			contentShort,
			createdAt,
			updatedAt,
		)
		pages = append(pages, *page)
	}

	if err = rows.Err(); err != nil {
		r.log.Debug().Err(err).Msg("ListPages2")
		return nil, nil, fmt.Errorf("error after iterating upload rows: %w", err)
	}

	// Totals
	whereClause, whereParams = sql.BuildWhere(params, "AND", 1)
	query = strings.Replace(queryPagesSelectTotal, "%WHERE%", "WHERE "+whereClause, 1)

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

	return &pages, &pagination, nil
}
