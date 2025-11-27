package handler

import (
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rafaeldepontes/auth-go/internal/domain"
)

// Handler controls the system routes based on *chi.Mux and a configuration struct.
func Handler(r *chi.Mux, app *domain.Application) {
	r.Use(chimiddleware.StripSlashes)

	// Public
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/users", app.UserService.FindAllUsers)
	})

	// Protected
	// WIP
}
