package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/rafaeldepontes/fauthless-go/internal/cache"
	"github.com/rafaeldepontes/fauthless-go/internal/errorhandler"
	jwt "github.com/rafaeldepontes/fauthless-go/internal/token"
)

type Middleware struct {
	JwtBuilder *jwt.JwtBuilder
	UserCache  *cache.Cache[string, string]
	Cache      *cache.Caches
}

type contextKey string

const TokenContextKey = contextKey("token")

var Token_Prefix = "Bearer "

func NewMiddleware(sk string, cache *cache.Caches) *Middleware {
	return &Middleware{
		JwtBuilder: jwt.NewJwtBuilder(sk),
		Cache:      cache,
	}
}

func (m *Middleware) AuthCookieBased(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var sessioToken *http.Cookie
		sessioToken, err := r.Cookie("session_token")

		userCache := m.Cache.UserCache
		if _, ok := userCache.Get("session_token"); !ok {
			errorhandler.UnauthroizedErrorHandler(w, errorhandler.ErrInvalidToken)
			return
		}

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

		if dirtToken == "" || !strings.HasPrefix(dirtToken, Token_Prefix) {
			fmt.Println(dirtToken == "" || !strings.HasPrefix(dirtToken, Token_Prefix))
			errorhandler.UnauthroizedErrorHandler(w, errorhandler.ErrInvalidToken)
			return
		}

		userClaims, err := m.JwtBuilder.VerifyToken(token)
		if err != nil {
			fmt.Println(err)
			errorhandler.UnauthroizedErrorHandler(w, err)
			return
		}

		isRefresh := false
		if !validCredentials(w, r, m, &token, userClaims, isRefresh) {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) JwtRefreshBased(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dirtToken := r.Header.Get("Authorization")
		token := cleanToken(dirtToken)

		if dirtToken == "" || !strings.HasPrefix(dirtToken, Token_Prefix) {
			fmt.Println(dirtToken == "" || !strings.HasPrefix(dirtToken, Token_Prefix))
			errorhandler.UnauthroizedErrorHandler(w, errorhandler.ErrInvalidToken)
			return
		}

		userClaims, err := m.JwtBuilder.VerifyToken(token)
		if err != nil {
			fmt.Println(err)
			errorhandler.UnauthroizedErrorHandler(w, err)
			return
		}

		isRefresh := true
		if !validCredentials(w, r, m, &token, userClaims, isRefresh) {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func cleanToken(dirtToken string) string {
	return strings.TrimPrefix(dirtToken, Token_Prefix)
}

func validCredentials(w http.ResponseWriter, r *http.Request, m *Middleware, token *string, userClaims *jwt.UserClaims, isRefresh bool) bool {
	path := r.URL.Path

	if !strings.Contains(path, "users") {
		return true
	}

	if !isRefresh {
		tokenCache := m.Cache.TokenCache
		if val, ok := tokenCache.Get(*token); !ok || val {
			fmt.Println("Missing cache")
			errorhandler.UnauthroizedErrorHandler(w, errorhandler.ErrInvalidToken)
			return false
		}
	}

	if !checkMethods(r) {
		return true
	}

	if r.Method == http.MethodDelete {
		tokenCache := m.Cache.TokenCache
		tokenCache.Set(*token, true, userClaims.ExpiresAt.Time)
	}

	pathSlice := strings.Split(path, "/")
	if len(pathSlice) <= 0 {
		fmt.Println("Invalid path")
		errorhandler.BadRequestErrorHandler(w, errorhandler.ErrIdIsRequired, r.URL.Path)
		return false
	}

	// for some reason when I try to get the path id by the normal way
	// using r.PathValue("username"), it's not giving me anything even though
	// I have the id in the request and have confirmed it... so this
	// was the only way to work around...
	username := pathSlice[len(pathSlice)-1]

	if userClaims.Username == username {
		return true
	}

	errorhandler.ForbiddenErrorHandler(w, errorhandler.ErrInvalidId)
	return false
}

func checkMethods(r *http.Request) bool {
	return r.Method == http.MethodPatch || r.Method == http.MethodDelete
}
