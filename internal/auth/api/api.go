package router

import (
	"net/http"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc(`GET /user`, http.HandlerFunc(HandlerUser))
	return mux
}
