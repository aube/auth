package page

import (
	"context"
	"strconv"
)

func (s *PageService) DeleteForce(ctx context.Context, id int64) error {
	err := s.repo.DeleteForce(ctx, id)
	if err != nil {
		s.log.Debug().Err(err).Msg("Delete")
		return err
	}
	s.log.Debug().Msg("DELETE! page: " + strconv.Itoa(int(id)))
	return nil
}
