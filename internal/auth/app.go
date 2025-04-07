package authserver

import (
	"fmt"
	"net/http"

	"github.com/aube/gophermart/internal/auth/api"
	"github.com/aube/gophermart/internal/auth/store"
)

/* type ApiServer interface {
	// logger       *logrus.Logger
	// store  store.Store
	router http.ServeMux
}
*/

// Start ...
func Start() error {
	config := NewConfig()

	fmt.Println("config.DatabaseDSN", config.DatabaseDSN)

	store, err := store.NewStore(config.DatabaseDSN)

	if err != nil {
		panic(err)
	}

	// defer store.db.Close()

	mux := api.NewRouter(store)

	return http.ListenAndServe(config.ServerAddress, mux)
}
