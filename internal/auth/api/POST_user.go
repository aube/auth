package api

import (
	"io"
	"net/http"

	"github.com/aube/gophermart/internal/auth/model"
)

func (s *Server) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	httpStatus := http.StatusCreated

	ctx := r.Context()

	if r.Body == nil || r.ContentLength == 0 {
		s.logger.ErrorContext(ctx, "HandlerCreateUser", "Request body is empty", "")
		http.Error(w, "Request body is empty", http.StatusBadRequest)
		return
	}

	// Body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		s.logger.ErrorContext(ctx, "HandlerCreateUser", "err", err)
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// JSON
	user, err := model.ParseCredentials(body)
	if err != nil {
		s.logger.ErrorContext(ctx, "HandlerCreateUser", "err", err)
		return
	}

	// Store
	err = s.store.User().Create(ctx, &user)
	if err != nil {
		s.logger.ErrorContext(ctx, "HandlerCreateUser", "err", err)
		httpStatus = http.StatusConflict
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	w.Write([]byte("Ololo, World!"))

	s.logger.Debug("HandlerCreateUser", "httpStatus", err)
}
