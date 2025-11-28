package api

import (
	"github.com/rafaeldepontes/auth-go/internal/service"
	log "github.com/sirupsen/logrus"
)

type Application struct {
	UserService *service.UserService
	AuthService *service.AuthService
	Logger      *log.Logger
}
