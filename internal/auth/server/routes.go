package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/rafaeldepontes/auth-go/internal/auth"
)

func MapAuthRoutes(r *chi.Mux, controller *auth.Controller) {
	(*r).Post("/register", (*controller).RegisterEp)
}

func MapAuthRoutesCookie(r *chi.Mux, controller *auth.Controller) {
	(*r).Post("/login", (*controller).LoginCookieBasedEp)
}

func MapAuthRoutesJwt(r *chi.Mux, controller *auth.Controller) {
	(*r).Post("/login", (*controller).LoginJwtBasedEp)
}

func MapAuthRoutesJwtRefresh(r *chi.Mux, controller *auth.Controller) {
	(*r).Post("/login", (*controller).LoginJwtRefreshBasedEp)
	(*r).Post("/renew", (*controller).RenewAccessTokenEp)
	(*r).Patch("/revoke/{id}", (*controller).RevokeSessionEp)
}

func MapAuthRoutesOAuth2(r *chi.Mux, controller *auth.Controller) {
	(*r).Get("/auth/{prodiver}/callback", (*controller).GetAuthCallbackOAuth2Ep)
	(*r).Get("/logout/{provider}", (*controller).LogoutOAuth2Ep)
	(*r).Get("/auth/{provider}", (*controller).GetAuthOAuth2Ep)
}
