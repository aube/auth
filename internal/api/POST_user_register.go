package api

import (
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/aube/auth/internal/helpers"
	"github.com/aube/auth/internal/httperrors"
	"github.com/aube/auth/internal/model"
)

func NewUserRegisterHandler(
	storeUser UserProvider,
	storeActiveUser ActiveUserProvider,
	logger *slog.Logger,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		if r.Body == nil || r.ContentLength == 0 {
			logger.ErrorContext(ctx, "HandlerCreateUser", "Request body is empty", "")
			http.Error(w, "Request body is empty", http.StatusBadRequest)
			return
		}

		// Body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.ErrorContext(ctx, "HandlerCreateUser", "err", err)
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		// JSON
		user, err := model.ParseCredentials(body)
		if err != nil {
			logger.ErrorContext(ctx, "HandlerCreateUser", "err", err)
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}

		// Store
		err = storeUser.Register(ctx, &user)
		if err != nil {
			logger.ErrorContext(ctx, "HandlerCreateUser", "err", err)

			var heherr *httperrors.HTTPError
			if errors.As(err, &heherr) {
				http.Error(w, heherr.Message, heherr.Code)
			} else {
				http.Error(w, "Failed to create user", http.StatusInternalServerError)
			}

			return
		}

		user.AfterLogin()

		storeActiveUser.Set(ctx, &user)

		helpers.SetAuthCookie(w, user.RandomHash)

		w.Header().Set("Authorization", bearerString+user.RandomHash)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("User registered"))

		logger.Debug("HandlerCreateUser", "User registered", user.ID)
	}
}
