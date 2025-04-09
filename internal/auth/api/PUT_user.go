package api

import (
	"net/http"
)

func HandlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}
