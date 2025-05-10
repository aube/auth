package api

import (
	"log/slog"
	"net/http"

	"github.com/aube/auth/internal/logger"
	"github.com/aube/auth/internal/store"
)

type Server struct {
	logger *slog.Logger
	router *http.ServeMux
	store  store.Store
}

func NewRouter(
	storeActiveUser ActiveUserProvider,
	storeUser UserProvider,
) *http.ServeMux {
	mux := http.NewServeMux()
	logger := logger.New()
	AuthMiddleware := NewAuthMiddleware(storeActiveUser, logger)

	// Public
	mux.HandleFunc(`POST /api/user/register`, NewUserRegisterHandler(storeUser, storeActiveUser, logger))
	mux.HandleFunc(`POST /api/user/login`, NewUserLoginHandler(storeUser, storeActiveUser, logger))

	// Private
	mux.HandleFunc(`GET /api/user`, AuthMiddleware(NewUserHanlder(storeUser, storeActiveUser, logger)))
	mux.HandleFunc(`DELETE /api/user`, AuthMiddleware(NewUserDeleteHandler(storeUser, storeActiveUser, logger)))

	return mux
}
