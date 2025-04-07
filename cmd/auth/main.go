package main

import (
	"fmt"

	authserver "github.com/aube/gophermart/internal/auth"
)

func main() {
	fmt.Println("gophermart auth running")

	authserver.Start()
}
