package middleware

import "net/http"

func AuthCookieBased(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		


		next.ServeHTTP(w, r)
	})
}

func JwtBased(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		


		next.ServeHTTP(w, r)
	})
}

func JwtRefreshBased(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		


		next.ServeHTTP(w, r)
	})
}
