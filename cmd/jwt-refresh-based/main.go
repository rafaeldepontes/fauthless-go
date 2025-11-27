package main

import (
	"fmt"
	"net/http"

	"github.com/rafaeldepontes/auth-go/api"
)

func main() {
	config := api.Init()

	fmt.Printf("API running at %v\n", config.JwtRefreshBasedPort)

	http.ListenAndServe(config.JwtRefreshBasedPort, nil)
}
