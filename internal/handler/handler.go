package handler

import (
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rafaeldepontes/auth-go/api"
)

// Handler controls the system routes based on *chi.Mux and a configuration struct.
func Handler(r *chi.Mux, app *api.Application, typeOf int) {
	r.Use(chimiddleware.StripSlashes)

	// Public
	switch typeOf {
	case api.CookieBased:
		r.Post("/login", app.AuthService.LoginCookieBased)
	case api.JwtBased:
		r.Post("/login", app.AuthService.LoginJwtBased)
	case api.JwtRefreshBased:
		r.Post("/login", app.AuthService.LoginJwtRefreshBased)
		r.Post("/renew", app.AuthService.RenewAccessToken)
		r.Post("/revoke", app.AuthService.RevokeSession)
	default:
		app.Logger.Fatalln("No authentication method was chosen.")
	}
	r.Post("/register", app.AuthService.Register)

	// Protected
	r.Group(func(r chi.Router) {
		r.Route("/api/v1", func(r chi.Router) {
			switch typeOf {
			case api.CookieBased:
				r.Use(app.Middleware.AuthCookieBased)
			case api.JwtBased:
				r.Use(app.Middleware.JwtBased)
			case api.JwtRefreshBased:
				r.Use(app.Middleware.JwtRefreshBased) // TODO: FINISH THE IMPLEMENTATION...
			default:
				app.Logger.Fatalln("No authentication method was chosen.")
			}

			r.Get("/users", app.UserService.FindAllUsers)
			r.Get("/users/{id}", app.UserService.FindUserById)

			switch typeOf {
			case api.CookieBased:
				app.Logger.Infoln("Cookie based authorization doens't allow this endpoints...")
			default:
				r.Patch("/users/{username}", app.UserService.UpdateUserDetails)
				r.Delete("/users/{username}", app.UserService.DeleteAccount)
			}
		})
	})
}
