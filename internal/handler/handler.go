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
		r.Post("/revoke/{id}", app.AuthService.RevokeSession)
	case api.OAuth2:
		r.Get("/auth/{prodiver}/callback", app.AuthService.GetAuthCallbackOAuth2)
		r.Get("/logout/{provider}", app.AuthService.LogoutOAuth2)
		r.Get("/auth/{provider}", app.AuthService.GetAuthOAuth2)
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
				r.Use(app.Middleware.JwtRefreshBased)
			case api.OAuth2:

			default:
				app.Logger.Fatal("No authentication method was chosen.")
			}

			// I could have done this in the same request, but for the learning purposes,
			// I'm doing it separately.
			r.Get("/users/cursor-pagination", app.UserService.FindAllUsersCursorPagination)
			r.Get("/users/offset-pagination", app.UserService.FindAllUsersOffSetPagination)
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
