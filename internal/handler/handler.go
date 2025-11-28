package handler

import (
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rafaeldepontes/auth-go/api"
	"github.com/rafaeldepontes/auth-go/internal/middleware"
)

// Handler controls the system routes based on *chi.Mux and a configuration struct.
func Handler(r *chi.Mux, app *api.Application, typeOf int) {
	r.Use(chimiddleware.StripSlashes)

	// Public
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/login", app.AuthService.Login)
		r.Get("/register", app.AuthService.Register)
	})

	switch typeOf {
	case api.CookieBased:
		r.Group(func(r chi.Router) {
			r.Route("/api/v1", func(r chi.Router) {
				r.Use(middleware.AuthCookieBased) // TODO: FINISH THE IMPLEMENTATION...
				r.Get("/users", app.UserService.FindAllUsers)
				r.Get("/users/{id}", app.UserService.FindUserById)
				r.Patch("/users", app.UserService.UpdateUserDetails)
				r.Delete("/users", app.UserService.DeleteAccount)
			})
		})
	case api.JwtBased:
		r.Group(func(r chi.Router) {
			r.Route("/api/v1", func(r chi.Router) {
				r.Use(middleware.JwtBased) // TODO: FINISH THE IMPLEMENTATION...
				r.Get("/users", app.UserService.FindAllUsers)
				r.Get("/users/{id}", app.UserService.FindUserById)
				r.Patch("/users", app.UserService.UpdateUserDetails)
				r.Delete("/users", app.UserService.DeleteAccount)
			})
		})
	case api.JwtRefreshBased:
		r.Group(func(r chi.Router) {
			r.Route("/api/v1", func(r chi.Router) {
				r.Use(middleware.JwtRefreshBased) // TODO: FINISH THE IMPLEMENTATION...
				r.Get("/users", app.UserService.FindAllUsers)
				r.Get("/users/{id}", app.UserService.FindUserById)
				r.Patch("/users", app.UserService.UpdateUserDetails)
				r.Delete("/users", app.UserService.DeleteAccount)
			})
		})
	default:
		app.Logger.Fatalln("No authentication method was chosen.")
	}
}
