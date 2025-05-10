package authserver

import (
	"net/http"

	"github.com/aube/auth/internal/api"
	"github.com/aube/auth/internal/store"
)

// Start ...
func Start() error {
	config := NewConfig()

	store, err := store.NewStore(config.DatabaseDSN)

	if err != nil {
		panic(err)
	}

	// defer store.db.Close()

	mux := api.NewRouter(store)

	return http.ListenAndServe(config.ServerAddress, mux)
}
