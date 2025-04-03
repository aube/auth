package router

import (
	"net/http"
)

func HandlerDeleteUserServiceAccess(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
