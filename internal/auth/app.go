package authserver

import (
	"net/http"

	"github.com/aube/gophermart/internal/auth/router"
	"github.com/aube/gophermart/internal/auth/store"
)

// Start ...
func Start(config *Config) error {
	_, err := store.NewStore(config.DatabaseURL)

	if err != nil {
		return err
	}

	// defer store.db.Close()

	router := router.NewRouter()

	return http.ListenAndServe(config.BindAddr, router)
}
