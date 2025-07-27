package dto_test

import (
	"testing"

	"github.com/aube/auth/internal/application/dto"
	"github.com/stretchr/testify/assert"
)

func TestPagination_Defaults(t *testing.T) {
	p := dto.Pagination{}
	assert.Equal(t, 0, p.Size)
	assert.Equal(t, 0, p.Page)
	assert.Equal(t, 0, p.Total)
}

func TestPagination_Fields(t *testing.T) {
	p := dto.Pagination{
		Size:  10,
		Page:  2,
		Total: 100,
	}

	assert.Equal(t, 10, p.Size)
	assert.Equal(t, 2, p.Page)
	assert.Equal(t, 100, p.Total)
}
