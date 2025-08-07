package page

import (
	"context"
	"errors"

	"github.com/aube/auth/internal/application/dto"
	"github.com/aube/auth/internal/domain/entities"
)

var ErrPageNotFound = errors.New("page not found")

type PageRepository interface {
	Create(ctx context.Context, page *entities.Page) error
	Update(ctx context.Context, page *entities.Page) error
	GetIDByName(ctx context.Context, name string) (int64, error)
	Delete(ctx context.Context, id int64) error
	DeleteForce(ctx context.Context, id int64) error

	FindByName(ctx context.Context, name string) (*entities.PageWithTime, error)
	FindByID(ctx context.Context, id int64) (*entities.PageWithTime, error)
	ListPages(ctx context.Context, limit int, offset int, params map[string]any) (*entities.PagesWithTimes, *dto.Pagination, error)
}
