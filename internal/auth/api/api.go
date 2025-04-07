package api

import (
	"net/http"

	"github.com/aube/gophermart/internal/auth/store"
)

type server struct {
	// logger       *logrus.Logger
	store  store.Store
	router *http.ServeMux
}

func NewRouter(store store.Store) *http.ServeMux {
	mux := http.NewServeMux()

	s := &server{
		store:  store,
		router: mux,
	}

	s.configureRouter()

	return s.router
}

func (s *server) configureRouter() {
	s.router.HandleFunc(`GET /user`, http.HandlerFunc(s.HandlerUser))
}
