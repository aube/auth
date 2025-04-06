package api

import (
	"net/http"
)

func HandlerUser(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("x-token")

	if token == "" {
		http.Error(w, "x-token header must be specified", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
