// Package entities defines the core domain models for the application.
package entities

// File represents a physical file in the system.
// Fields:
//   - Name: Unique identifier for the file (typically a UUID)
//   - Size: File size in bytes
//   - Path: Storage location path (filesystem, S3, etc.)
// JSON tags support serialization for API responses.

type File struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
	Path string `json:"path"`
}

// Files is a collection type for multiple File entities.
type Files []File

// NewFile creates a new File instance.
// name: Unique file identifier
// path: Storage system path
// size: File size in bytes
// Returns: *File instance
// Note: Name typically represents the server-side UUID filename
func NewFile(name, path string, size int64) *File {
	return &File{
		Name: name,
		Size: size,
		Path: path,
	}
}
