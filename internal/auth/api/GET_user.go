package api

import (
	"net/http"
)

func (s *server) HandlerUser(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("x-token")

	s.logger.Info("ololo")

	if token == "" {
		http.Error(w, "x-token header must be specified", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
