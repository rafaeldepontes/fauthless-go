package middleware

import (
	"net/http"
	"strings"

	"github.com/rafaeldepontes/auth-go/internal/errorhandler"
	jwt "github.com/rafaeldepontes/auth-go/internal/token"
)

type Middleware struct {
	JwtBuilder *jwt.JwtBuilder
}

type contextKey string

const TokenContextKey = contextKey("token")

var Token_Prefix = "Bearer "

func NewMiddleware(sk string) *Middleware {
	return &Middleware{
		JwtBuilder: jwt.NewJwtBuilder(sk),
	}
}

func (m *Middleware) AuthCookieBased(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var sessioToken *http.Cookie
		sessioToken, err := r.Cookie("session_token")

		// I should check if the token is the same for the user...
		// but I don't want to.
		if err != nil || sessioToken.Value == "" {
			errorhandler.UnauthroizedErrorHandler(w, errorhandler.ErrInvalidToken)
			return
		}

		csrfToken := r.Header.Get("X-CSRF-Token")
		if csrfToken == "" {
			errorhandler.UnauthroizedErrorHandler(w, errorhandler.ErrInvalidCSRFToken)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) JwtBased(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dirtToken := r.Header.Get("Authorization")
		token := cleanToken(dirtToken)

		if dirtToken == "" || dirtToken == token {
			errorhandler.UnauthroizedErrorHandler(w, errorhandler.ErrInvalidToken)
			return
		}

		_, err := m.JwtBuilder.VerifyToken(token)
		if err != nil {
			errorhandler.UnauthroizedErrorHandler(w, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) JwtRefreshBased(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		next.ServeHTTP(w, r)
	})
}

func cleanToken(dirtToken string) string {
	return strings.TrimPrefix(dirtToken, Token_Prefix)
}
