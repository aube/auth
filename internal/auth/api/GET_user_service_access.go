package api

import (
	"net/http"
)

func HandlerUserServiceAccess(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
