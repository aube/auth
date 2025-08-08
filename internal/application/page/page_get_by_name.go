package page

import (
	"context"

	"github.com/aube/auth/internal/domain/entities"
)

func (s *PageService) GetByName(ctx context.Context, name string) (*entities.PageWithTime, error) {
	page, err := s.repo.FindByName(ctx, name)
	if err != nil {
		s.log.Debug().Err(err).Msg("GetByName")
		return nil, err
	}

	return page, nil
}
