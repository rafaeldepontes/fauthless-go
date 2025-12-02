package auth

import (
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"

	"github.com/rafaeldepontes/auth-go/configs"
)

const (
	MaxAge = 86400 * 30
	IsProd = false
)

func InitOAuth(config *configs.Configuration) {
	var store *sessions.CookieStore = sessions.NewCookieStore([]byte(config.GoogleSecretKey))
	store.MaxAge(MaxAge)

	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = IsProd

	gothic.Store = store

	goth.UseProviders(google.New(config.GoogleSecretKey, config.GoogleClientSecret, config.UrlCallback))
}
