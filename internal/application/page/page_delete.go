package page

import (
	"context"
	"strconv"
)

func (s *PageService) Delete(ctx context.Context, id int64) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.log.Debug().Err(err).Msg("Delete")
		return err
	}

	s.log.Debug().Msg("DELETE page: " + strconv.Itoa(int(id)))
	return nil
}
