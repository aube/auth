// Package dto contains data transfer objects for API requests and responses.
package dto

// Pagination represents pagination metadata for API responses.
// Fields:
//   - Size: Number of items per page (default: 10).
//   - Page: Current page number (default: 1).
//   - Total: Total number of items available.
//
// Used in list operations to provide pagination context.
type Pagination struct {
	Size  int `json:"size"`
	Page  int `json:"page"`
	Total int `json:"total"`
}
