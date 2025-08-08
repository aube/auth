package page

import (
	"context"
	"errors"
	"strconv"

	"github.com/aube/auth/internal/application/dto"
	"github.com/aube/auth/internal/domain/entities"
)

func (s *PageService) Create(ctx context.Context, pageDTO dto.CreatePageRequest) (*entities.PageWithTime, error) {
	// Проверяем, существует ли страница с таким именем
	id, err := s.repo.GetIDByName(ctx, pageDTO.Name)
	if err != nil {
		s.log.Debug().Err(err).Msg("Create1")
		return nil, err
	}
	if id > 0 {
		return nil, errors.New("page already exists")
	}

	// Создаем сущность page
	page, err := entities.NewPage(
		0,
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
		s.log.Debug().Err(err).Msg("Create2")
		return nil, err
	}

	// Сохраняем в репозитории
	if err := s.repo.Create(ctx, page); err != nil {
		s.log.Debug().Err(err).Msg("Create3")
		return nil, err
	}

	// Получаем сохранённый результат
	createdPage, err := s.repo.FindByID(ctx, page.ID)
	if err != nil {
		s.log.Debug().Err(err).Msg("Update4")
		return nil, err
	}

	s.log.Debug().Msg("CREATE page: " + strconv.Itoa(int(page.ID)) + ", " + pageDTO.Name)
	return createdPage, nil
}
