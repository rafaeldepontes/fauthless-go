package api

import (
	"github.com/rafaeldepontes/fauthless-go/internal/auth"
	"github.com/rafaeldepontes/fauthless-go/internal/middleware"
	"github.com/rafaeldepontes/fauthless-go/internal/user"
	log "github.com/sirupsen/logrus"
)

type Application struct {
	UserController *user.Controller
	AuthController *auth.Controller
	Logger         *log.Logger
	Middleware     *middleware.Middleware
}
