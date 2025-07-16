package entities

type Upload struct {
	file        File
	ID          string `json:"id"`
	Description string `json:"description"`
}

func NewUpload(file File, id, description string) *Upload {
	return &Upload{
		file:        file,
		ID:          id,
		Description: description,
	}
}
