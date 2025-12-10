package server

import (
	"net/http"

	"github.com/rafaeldepontes/fauthless-go/internal/auth"
)

type authController struct {
	service *auth.Service
}

func NewAuthController(s *auth.Service) auth.Controller {
	return &authController{
		service: s,
	}
}

func (s *authController) RegisterEp(w http.ResponseWriter, r *http.Request) {
	(*s.service).Register(w, r)
}

func (s *authController) LoginCookieBasedEp(w http.ResponseWriter, r *http.Request) {
	(*s.service).LoginCookieBased(w, r)
}

func (s *authController) LoginJwtBasedEp(w http.ResponseWriter, r *http.Request) {
	(*s.service).LoginJwtBased(w, r)
}

func (s *authController) LoginJwtRefreshBasedEp(w http.ResponseWriter, r *http.Request) {
	(*s.service).LoginJwtRefreshBased(w, r)
}

func (s *authController) RenewAccessTokenEp(w http.ResponseWriter, r *http.Request) {
	(*s.service).RenewAccessToken(w, r)
}

func (s *authController) RevokeSessionEp(w http.ResponseWriter, r *http.Request) {
	(*s.service).RevokeSession(w, r)
}

func (s *authController) GetAuthCallbackOAuth2Ep(w http.ResponseWriter, r *http.Request) {
	(*s.service).GetAuthCallbackOAuth2(w, r)
}

func (s *authController) LogoutOAuth2Ep(w http.ResponseWriter, r *http.Request) {
	(*s.service).LogoutOAuth2(w, r)
}

func (s *authController) GetAuthOAuth2Ep(w http.ResponseWriter, r *http.Request) {
	(*s.service).GetAuthOAuth2(w, r)
}
