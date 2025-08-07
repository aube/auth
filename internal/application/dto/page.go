package dto

import (
	"time"

	"github.com/aube/auth/internal/domain/entities"
)

type CreatePageRequest struct {
	Name         string `json:"name"`
	Meta         string `json:"meta"`
	Title        string `json:"title"`
	Category     string `json:"category"`
	Template     string `json:"template"`
	H1           string `json:"h1"`
	Content      string `json:"content"`
	ContentShort string `json:"content_short"`
}

type UpdatePageRequest struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Meta         string `json:"meta"`
	Title        string `json:"title"`
	Category     string `json:"category"`
	Template     string `json:"template"`
	H1           string `json:"h1"`
	Content      string `json:"content"`
	ContentShort string `json:"content_short"`
}

type PageResponse struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Meta         string    `json:"meta"`
	Title        string    `json:"title"`
	Category     string    `json:"category"`
	Template     string    `json:"template"`
	H1           string    `json:"h1"`
	Content      string    `json:"content"`
	ContentShort string    `json:"content_short"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func NewPageResponse(page *entities.PageWithTime) *PageResponse {
	return &PageResponse{
		ID:           page.ID,
		Name:         page.Name,
		Meta:         page.Meta,
		Title:        page.Title,
		Category:     page.Category,
		Template:     page.Template,
		H1:           page.H1,
		Content:      page.Content,
		ContentShort: page.ContentShort,
		CreatedAt:    page.CreatedAt,
		UpdatedAt:    page.UpdatedAt,
	}
}
