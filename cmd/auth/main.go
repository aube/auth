package main

import (
	"fmt"

	authserver "github.com/aube/gophermart/internal/auth"
)

func main() {
	fmt.Println("gophermert auth")

	authserver.Start()
}
