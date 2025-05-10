package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/aube/auth/internal/ctxkeys"
)

const authCookieName = "auth"
const bearerString = "Bearer "

type Middleware func(http.HandlerFunc) http.HandlerFunc

func NewAuthMiddleware(
	storeActiveUser ActiveUserProvider,
	logger *slog.Logger,
) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var token string

			authHeader := r.Header.Get("Authorization")

			if authHeader != "" {
				token = authHeader[len(bearerString):]
			}

			if token == "" {
				cookie, err := r.Cookie(authCookieName)
				if err == nil {
					token = cookie.Value
				}
			}

			w.Header().Set("Authorization", bearerString+token)

			user, ok := storeActiveUser.Get(r.Context(), token)

			if !ok {
				http.Error(w, "Auth failed", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ctxkeys.UserID, user.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
