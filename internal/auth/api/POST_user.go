package api

import (
	"net/http"
)

func HandlerCreateUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}
