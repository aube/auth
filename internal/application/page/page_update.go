package page

import (
	"context"
	"errors"
	"strconv"

	"github.com/aube/auth/internal/application/dto"
	"github.com/aube/auth/internal/domain/entities"
)

func (s *PageService) Update(ctx context.Context, pageDTO dto.UpdatePageRequest) (*entities.PageWithTime, error) {

	// Проверяем, существует ли другая страница с таким именем
	id, err := s.repo.GetIDByName(ctx, pageDTO.Name)
	if err != nil {
		s.log.Debug().Err(err).Msg("Update1")
		return nil, err
	}
	if id > 0 && id != pageDTO.ID {
		return nil, errors.New("page with name " + pageDTO.Name + " already exists")
	}

	// Создаем сущность page
	page, err := entities.NewPage(
		pageDTO.ID,
		pageDTO.Name,
		pageDTO.Meta,
		pageDTO.Title,
		pageDTO.Category,
		pageDTO.Template,
		pageDTO.H1,
		pageDTO.Content,
		pageDTO.ContentShort,
	)
	if err != nil {
		s.log.Debug().Err(err).Msg("Update2")
		return nil, err
	}
	// Сохраняем в репозитории
	if err := s.repo.Update(ctx, page); err != nil {
		s.log.Debug().Err(err).Msg("Update3")
		return nil, err
	}
	// Получаем сохранённый результат
	updatedPage, err := s.repo.FindByID(ctx, pageDTO.ID)
	if err != nil {
		s.log.Debug().Err(err).Msg("Update4")
		return nil, err
	}

	s.log.Debug().Msg("UPDATE page: " + strconv.Itoa(int(pageDTO.ID)) + ", " + pageDTO.Name)
	return updatedPage, nil
}
