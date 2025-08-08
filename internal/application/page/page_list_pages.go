package page

import (
	"context"

	"github.com/aube/auth/internal/application/dto"
	"github.com/aube/auth/internal/domain/entities"
)

func (s *PageService) ListPages(ctx context.Context, offset, limit int, params map[string]any) (*entities.PagesWithTimes, *dto.Pagination, error) {
	return s.repo.ListPages(ctx, offset, limit, params)
}
