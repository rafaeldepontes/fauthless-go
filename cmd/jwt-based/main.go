package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rafaeldepontes/auth-go/api"
	"github.com/rafaeldepontes/auth-go/configs"
	"github.com/rafaeldepontes/auth-go/internal/domain"
	"github.com/rafaeldepontes/auth-go/internal/handler"
)

func main() {
	var app *domain.Application
	var config *configs.Configuration

	config, app, db, err := api.Init()
	if err != nil {
		app.Logger.Fatalf("[ERROR] An error occurred: %v", err)
	}
	defer db.Close()

	var r *chi.Mux = chi.NewRouter()
	handler.Handler(r, app)

	app.Logger.Infof("[LOG] API running at %v\n", config.JwtBasedPort)

	http.ListenAndServe(config.JwtBasedPort, nil)
}
