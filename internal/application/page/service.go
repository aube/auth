package page

import (
	"context"
	"errors"
	"strconv"

	"github.com/aube/auth/internal/application/dto"
	"github.com/aube/auth/internal/domain/entities"
	"github.com/aube/auth/internal/utils/logger"
	"github.com/rs/zerolog"
)

type PageService struct {
	repo PageRepository
	log  zerolog.Logger
}

func NewPageService(repo PageRepository) *PageService {
	return &PageService{
		repo: repo,
		log:  logger.Get().With().Str("page", "service").Logger(),
	}
}

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

func (s *PageService) GetByID(ctx context.Context, id int64) (*entities.PageWithTime, error) {
	page, err := s.repo.FindByID(ctx, id)
	if err != nil {
		s.log.Debug().Err(err).Msg("GetByID")
		return nil, err
	}

	return page, nil
}
func (s *PageService) GetByName(ctx context.Context, name string) (*entities.PageWithTime, error) {
	page, err := s.repo.FindByName(ctx, name)
	if err != nil {
		s.log.Debug().Err(err).Msg("GetByName")
		return nil, err
	}

	return page, nil
}
func (s *PageService) ListPages(ctx context.Context, offset, limit int, params map[string]any) (*entities.PagesWithTimes, *dto.Pagination, error) {
	return s.repo.ListPages(ctx, offset, limit, params)
}

func (s *PageService) Delete(ctx context.Context, id int64) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		s.log.Debug().Err(err).Msg("Delete")
		return err
	}

	s.log.Debug().Msg("DELETE page: " + strconv.Itoa(int(id)))
	return nil
}

func (s *PageService) DeleteForce(ctx context.Context, id int64) error {
	err := s.repo.DeleteForce(ctx, id)
	if err != nil {
		s.log.Debug().Err(err).Msg("Delete")
		return err
	}
	s.log.Debug().Msg("DELETE! page: " + strconv.Itoa(int(id)))
	return nil
}
