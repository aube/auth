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

type Server struct {
	router     *gin.Engine
	httpServer *http.Server
}

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

func (s *Server) Start() error {
	log.Println("Server starting on :8080")
	return s.httpServer.ListenAndServe()
}

func (s *Server) Close() error {
	return s.httpServer.Close()
}
