package main

import (
	"fmt"
	"net/http"

	"github.com/aube/go-mart/internal/auth/router"
)

func main() {
	portNumber := "8081"
	mux := router.NewRouter()
	err := http.ListenAndServe(":"+portNumber, mux)

	if err != nil {
		fmt.Println("Error starting server:", err)
	} else {
		fmt.Println("Server started at:", portNumber)
	}
}
