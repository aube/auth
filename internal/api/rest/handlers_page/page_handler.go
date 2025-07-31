package handlers_page

import (
	"context"
	"net/http"
	"strconv"

	"github.com/aube/auth/internal/application/dto"
	"github.com/aube/auth/internal/utils/logger"
	"github.com/rs/zerolog"

	"github.com/gin-gonic/gin"
)

type PageService interface {
	Create(ctx context.Context, pageDTO dto.CreatePageRequest) (*dto.PageResponse, error)
	Update(ctx context.Context, pageDTO dto.UpdatePageRequest) (*dto.PageResponse, error)
	GetByID(ctx context.Context, id int64) (*dto.PageResponse, error)
	GetByName(ctx context.Context, name string) (*dto.PageResponse, error)
	Delete(ctx context.Context, id int64) error
	DeleteForce(ctx context.Context, id int64) error
}

type PageHandler interface {
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	GetByParam(c *gin.Context)
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

	createdUser, err := h.pageService.Create(ctx, pageDTO)
	if err != nil {
		h.log.Debug().Err(err).Msg("Create2")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdUser)
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

	createdUser, err := h.pageService.Update(ctx, pageDTO)
	if err != nil {
		h.log.Debug().Err(err).Msg("Update2")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdUser)
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

	c.JSON(http.StatusOK, page)
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

	c.JSON(http.StatusOK, page)
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

	c.Status(http.StatusNoContent)
}
