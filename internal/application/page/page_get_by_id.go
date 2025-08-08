package page

import (
	"context"

	"github.com/aube/auth/internal/domain/entities"
)

func (s *PageService) GetByID(ctx context.Context, id int64) (*entities.PageWithTime, error) {
	page, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.log.Debug().Err(err).Msg("GetByID")
		return nil, err
	}

	return page, nil
}
