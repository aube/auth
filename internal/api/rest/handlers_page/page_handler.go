package handlers_page

import (
	"context"
	"net/http"
	"strconv"

	"github.com/aube/auth/internal/application/dto"
	"github.com/aube/auth/internal/domain/entities"
	"github.com/aube/auth/internal/utils/logger"
	"github.com/rs/zerolog"

	"github.com/gin-gonic/gin"
)

type PageService interface {
	Delete(ctx context.Context, id int64) error
	DeleteForce(ctx context.Context, id int64) error

	Create(ctx context.Context, pageDTO dto.CreatePageRequest) (*entities.PageWithTime, error)
	Update(ctx context.Context, pageDTO dto.UpdatePageRequest) (*entities.PageWithTime, error)
	GetByName(ctx context.Context, name string) (*entities.PageWithTime, error)
	GetByID(ctx context.Context, id int64) (*entities.PageWithTime, error)
	ListPages(ctx context.Context, offset int, limit int, params map[string]any) (*entities.PagesWithTimes, *dto.Pagination, error)
}

type PageHandler interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetByParam(c *gin.Context)
	ListPages(c *gin.Context)
}

type Handler struct {
	pageService PageService
	jwtSecret   []byte
	log         zerolog.Logger
}

func NewPageHandler(pageService PageService, jwtSecret string) PageHandler {
	return &Handler{
		pageService: pageService,
		jwtSecret:   []byte(jwtSecret),
		log:         logger.Get().With().Str("handlers", "page_handler").Logger(),
	}
}

func (h *Handler) Create(c *gin.Context) {
	var req dto.CreatePageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Debug().Err(err).Msg("Create1")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	pageDTO := dto.CreatePageRequest(req)

	page, err := h.pageService.Create(ctx, pageDTO)
	if err != nil {
		h.log.Debug().Err(err).Msg("Create2")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.NewPageResponse(page))
}

func (h *Handler) Update(c *gin.Context) {
	var req dto.UpdatePageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Debug().Err(err).Msg("Update1")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	pageDTO := dto.UpdatePageRequest(req)

	page, err := h.pageService.Update(ctx, pageDTO)
	if err != nil {
		h.log.Debug().Err(err).Msg("Update2")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.NewPageResponse(page))
}

func (h *Handler) GetByParam(c *gin.Context) {
	ID := c.Query("id")
	name := c.Query("name")
	if ID == "" && name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Page ID is required"})
		return
	}

	if ID == "" {
		h.GetByName(c)
	} else {
		h.GetByID(c)
	}
}

func (h *Handler) GetByID(c *gin.Context) {
	ID := c.Query("id")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Page ID is required"})
		return
	}

	pageID, err := strconv.Atoi(ID)
	if err != nil {
		h.log.Debug().Msg("Page ID is not ok")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized2"})
		return
	}

	ctx := c.Request.Context()
	page, err := h.pageService.GetByID(ctx, int64(pageID))
	if err != nil {
		h.log.Debug().Err(err).Msg("GetProfile")
		c.JSON(http.StatusNotFound, gin.H{"error": "page not found"})
		return
	}
	h.log.Debug().Msg(page.Name)

	c.JSON(http.StatusOK, dto.NewPageResponse(page))
}

func (h *Handler) GetByName(c *gin.Context) {
	name := c.Query("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Page ID is required"})
		return
	}
	ctx := c.Request.Context()
	page, err := h.pageService.GetByName(ctx, name)
	if err != nil {
		h.log.Debug().Err(err).Msg("GetProfile")
		c.JSON(http.StatusNotFound, gin.H{"error": "page not found"})
		return
	}
	h.log.Debug().Msg(page.Name)

	c.JSON(http.StatusOK, dto.NewPageResponse(page))
}

func (h *Handler) ListPages(c *gin.Context) {

	offset := c.GetInt("offset")
	limit := c.GetInt("limit")

	params := make(map[string]any)
	if c.Query("updated_at") != "" {
		params["updated_at >="] = c.Query("updated_at")
	}
	params["deleted"] = "false"

	pages, pagination, err := h.pageService.ListPages(c.Request.Context(), offset, limit, params)
	if err != nil {
		h.log.Debug().Err(err).Msg("ListFiles")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list files"})
		return
	}

	rows := make([]dto.PageResponse, len(*pages))
	for i, page := range *pages {
		rows[i] = *dto.NewPageResponse(&page)
	}

	c.JSON(http.StatusOK, gin.H{
		"rows":       rows,
		"pagination": pagination,
	})
}

func (h *Handler) Delete(c *gin.Context) {
	ID := c.Query("id")
	if ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Page ID is required"})
		return
	}

	pageID, err := strconv.Atoi(ID)
	if err != nil {
		h.log.Debug().Msg("Page ID is not ok")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized2"})
		return
	}

	force := c.Query("force")

	ctx := c.Request.Context()
	if force == "" {
		err = h.pageService.Delete(ctx, int64(pageID))
	} else {
		err = h.pageService.DeleteForce(ctx, int64(pageID))
	}
	if err != nil {
		h.log.Debug().Err(err).Msg("Delete")
		c.JSON(http.StatusNotFound, gin.H{"error": "page not found"})
		return
	}

	c.Status(http.StatusOK)
}
