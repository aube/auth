package api

import (
	"log/slog"
	"net/http"
)

func NewUserHanlder(
	storeUser UserProvider,
	storeActiveUser ActiveUserProvider,
	logger *slog.Logger,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
