package handler

import (
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rafaeldepontes/fauthless-go/api"
	authServer "github.com/rafaeldepontes/fauthless-go/internal/auth/server"
	"github.com/rafaeldepontes/fauthless-go/internal/user/server"
)

// Handler controls the system routes based on *chi.Mux and a configuration struct.
func Handler(r *chi.Mux, app *api.Application, typeOf int) {
	r.Use(chimiddleware.StripSlashes)

	// Public
	switch typeOf {
	case api.CookieBased:
		authServer.MapAuthRoutesCookie(r, app.AuthController)
	case api.JwtBased:
		authServer.MapAuthRoutesJwt(r, app.AuthController)
	case api.JwtRefreshBased:
		authServer.MapAuthRoutesJwtRefresh(r, app.AuthController)
	case api.OAuth2:
		authServer.MapAuthRoutesOAuth2(r, app.AuthController)
	default:
		app.Logger.Fatalln("No authentication method was chosen.")
	}
	authServer.MapAuthRoutes(r, app.AuthController)

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
			case api.OAuth2: //do nothing...
			default:
				app.Logger.Fatal("No authentication method was chosen.")
			}

			server.MapUserRoutes(&r, app.UserController)

			switch typeOf {
			case api.CookieBased:
				app.Logger.Infoln("Cookie based authorization doens't allow this endpoints...")
			default:
				server.MapUserRoutesJwt(&r, app.UserController)
			}
		})
	})
}
