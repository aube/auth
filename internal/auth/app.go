package authserver

import (
	"net/http"

	"github.com/aube/gophermart/internal/auth/api"
	"github.com/aube/gophermart/internal/auth/store"
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
