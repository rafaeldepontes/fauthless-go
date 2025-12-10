package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rafaeldepontes/fauthless-go/api"
	"github.com/rafaeldepontes/fauthless-go/configs"
	"github.com/rafaeldepontes/fauthless-go/internal/handler"
)

func main() {
	var app *api.Application
	var config *configs.Configuration

	config, app, db, err := api.Init()
	if err != nil {
		app.Logger.Fatalf("An error occurred: %v", err)
	}
	defer db.Close()

	var r *chi.Mux = chi.NewRouter()
	handler.Handler(r, app, api.CookieBased)

	app.Logger.Infof("API running at %v", config.CookieBasedPort)

	http.ListenAndServe(":"+config.CookieBasedPort, r)
}
