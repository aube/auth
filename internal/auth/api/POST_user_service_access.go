package api

import (
	"net/http"
)

func HandlerCreateUserServiceAccess(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}
