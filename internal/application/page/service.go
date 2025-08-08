package page

import (
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
