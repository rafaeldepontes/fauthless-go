package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rafaeldepontes/auth-go/api"
	"github.com/rafaeldepontes/auth-go/internal/handler"
)

func main() {
	var app *api.Application
	var config *api.Configuration

	config, app, db, err := api.Init()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var r *chi.Mux = chi.NewRouter()
	handler.Handler(r, app)

	fmt.Printf("API running at %v\n", config.JwtRefreshBasedPort)

	http.ListenAndServe(config.JwtRefreshBasedPort, nil)
}
