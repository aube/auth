package router

import (
	"net/http"
)

func HandlerDeleteUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
