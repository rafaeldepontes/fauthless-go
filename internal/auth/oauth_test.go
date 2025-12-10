package auth

import (
	"testing"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/rafaeldepontes/fauthless-go/configs"
)

const (
	googleSecretKey    = "SecretKey"
	googleClientSecret = "GoogleClientSecret"
	urlCallback        = "http://localhost:8000/auth/google/callback"
)

// TestInitOAuth verifies InitOAuth sets the gothic store and
// registers the google provider.
func TestInitOAuth(t *testing.T) {
	oldStore := gothic.Store
	defer func() { gothic.Store = oldStore }()

	//given
	config := configs.Configuration{
		GoogleSecretKey:    googleSecretKey,
		GoogleClientSecret: googleClientSecret,
		UrlCallback:        urlCallback,
	}

	//when
	InitOAuth(&config)

	//then
	if gothic.Store == nil {
		t.Fatal("gothic.Store is nil after InitOAuth")
	}

	cs, ok := gothic.Store.(*sessions.CookieStore)
	if !ok {
		t.Fatalf("gothic.Store has unexpected type: %T", gothic.Store)
	}

	if cs.Options.Path != "/" {
		t.Fatalf("expected cookie path '/', got %v", cs.Options.Path)
	}

	if !cs.Options.HttpOnly {
		t.Fatal("expected cookie HttpOnly=true")
	}

	if cs.Options.Secure {
		t.Fatal("expected cookie Secure=false in test mode")
	}

	if _, err := goth.GetProvider("google"); err != nil {
		t.Fatalf("google provider not registered: %v", err)
	}
}
