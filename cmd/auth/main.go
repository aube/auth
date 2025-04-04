package main

import (
	"fmt"
	"net/http"

	"github.com/aube/gophermart/internal/auth/router"
)

func main() {
	fmt.Println("gophermert auth")

	portNumber := "8081"
	mux := router.NewRouter()
	err := http.ListenAndServe(":"+portNumber, mux)

	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
