package dto

type Pagination struct {
	Size  int `json:"size"`
	Page  int `json:"page"`
	Total int `json:"total"`
}
