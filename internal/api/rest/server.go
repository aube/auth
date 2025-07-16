package rest

import (
	"log"
	"net/http"

	appFile "github.com/aube/auth/internal/application/file"
	appUser "github.com/aube/auth/internal/application/user"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router     *gin.Engine
	httpServer *http.Server
}

func NewServer(
	userService *appUser.UserService,
	fileService *appFile.FileService,
	jwtSecret string,
	apiPath string,
) *Server {
	router, apiGroup := NewRouter(apiPath)
	SetupUserRouter(apiGroup, userService, jwtSecret)
	SetupFilesRouter(apiGroup, fileService, jwtSecret)
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
