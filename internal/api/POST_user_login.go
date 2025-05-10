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

func NewUserLoginHandler(
	storeUser UserProvider,
	storeActiveUser ActiveUserProvider,
	logger *slog.Logger,
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		if r.Body == nil || r.ContentLength == 0 {
			logger.ErrorContext(ctx, "UserLogin", "Request body is empty", "")
			http.Error(w, "Request body is empty", http.StatusBadRequest)
			return
		}

		// Body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.ErrorContext(ctx, "UserLogin", "err", err)
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		// JSON
		user, err := model.ParseCredentials(body)
		if err != nil {
			logger.ErrorContext(ctx, "UserLogin", "err", err)
			return
		}

		// Store
		_, err = storeUser.Login(ctx, &user)
		if err != nil {
			logger.ErrorContext(ctx, "UserLogin", "err", err)

			var heherr *httperrors.HTTPError
			if errors.As(err, &heherr) {
				http.Error(w, heherr.Message, heherr.Code)
			} else {
				http.Error(w, "Login failed", http.StatusInternalServerError)
			}

			return
		}

		user.AfterLogin()

		storeActiveUser.Set(ctx, &user)

		helpers.SetAuthCookie(w, user.RandomHash)

		w.Header().Set("Authorization", bearerString+user.RandomHash)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(user.RandomHash))

		logger.DebugContext(ctx, "UserLogin", "user", user)
	}
}
