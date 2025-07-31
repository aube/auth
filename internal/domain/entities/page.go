package entities

import (
	"errors"
)

type Page struct {
	ID           int64
	Name         string
	Meta         string
	Title        string
	Category     string
	Template     string
	H1           string
	Content      string
	ContentShort string
}

func NewPage(
	id int64,
	name string,
	meta string,
	title string,
	category string,
	template string,
	h1 string,
	content string,
	contentShort string,
) (*Page, error) {

	if name == "" {
		return nil, errors.New("name cannot be empty")
	}

	return &Page{
		ID:           id,
		Name:         name,
		Meta:         meta,
		Title:        title,
		Category:     category,
		Template:     template,
		H1:           h1,
		Content:      content,
		ContentShort: contentShort,
	}, nil
}
