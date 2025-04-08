package api

import (
	"log/slog"
	"net/http"

	"github.com/aube/gophermart/internal/auth/logger"
	"github.com/aube/gophermart/internal/auth/store"
)

type server struct {
	logger *slog.Logger
	router *http.ServeMux
	store  store.Store
}

func NewRouter(store store.Store) *http.ServeMux {
	mux := http.NewServeMux()

	s := &server{
		logger: logger.New(),
		store:  store,
		router: mux,
	}

	s.configureRouter()

	return s.router
}

func (s *server) configureRouter() {
	s.router.HandleFunc(`GET /user`, http.HandlerFunc(s.HandlerUser))
}
