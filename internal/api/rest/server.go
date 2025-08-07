// Package rest provides the HTTP REST server implementation for the application.
package rest

import (
	"log"
	"net/http"

	appFile "github.com/aube/auth/internal/application/file"
	appImage "github.com/aube/auth/internal/application/image"
	appPage "github.com/aube/auth/internal/application/page"
	appUpload "github.com/aube/auth/internal/application/upload"
	appUser "github.com/aube/auth/internal/application/user"
	"github.com/gin-gonic/gin"
)

// Server represents the HTTP server with router and HTTP server configurations.
// router: The Gin router handling all routes.
// httpServer: The underlying HTTP server.
type Server struct {
	router     *gin.Engine
	httpServer *http.Server
}

// NewServer initializes a new Server instance with configured routes and services.
// userService: Service for user operations.
// fileService: Service for file storage operations.
// uploadService: Service for upload metadata operations.
// jwtSecret: Secret key for JWT token generation and validation.
// apiPath: Base path for API routes (e.g., "/api").
// Returns: A configured *Server instance.
func NewServer(
	userService *appUser.UserService,
	pageService *appPage.PageService,
	fileService *appFile.FileService,
	imgFileService *appFile.FileService,
	uploadService *appUpload.UploadService,
	imageService *appImage.ImageService,
	jwtSecret string,
	apiPath string,
) *Server {
	router, apiGroup := NewRouter(apiPath)
	SetupUserRouter(apiGroup, userService, jwtSecret)
	SetupPageRouter(apiGroup, pageService, jwtSecret)
	SetupUploadsRouter(apiGroup, fileService, uploadService, jwtSecret)
	SetupImagesRouter(apiGroup, imgFileService, imageService, jwtSecret)
	SetupStaticRouter(router, apiPath)

	return &Server{
		router: router,
		httpServer: &http.Server{
			Addr:    ":8080",
			Handler: router,
		},
	}
}

// Start begins listening on the configured address (default ":8080").
// Returns: An error if the server fails to start.
func (s *Server) Start() error {
	log.Println("Server starting on :8080")
	return s.httpServer.ListenAndServe()
}

// Close shuts down the HTTP server gracefully.
// Returns: An error if the server fails to close.
func (s *Server) Close() error {
	return s.httpServer.Close()
}
