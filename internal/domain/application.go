package domain

import (
	"github.com/rafaeldepontes/auth-go/internal/service"
	log "github.com/sirupsen/logrus"
)

type Application struct {
	UserService *service.UserService
	Logger      *log.Logger
}
