package main

import (
	"fmt"

	authserver "github.com/aube/auth/internal/auth"
)

func main() {
	fmt.Println("auth auth running")

	authserver.Start()
}
