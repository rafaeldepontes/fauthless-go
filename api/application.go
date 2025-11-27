package api

import "github.com/rafaeldepontes/auth-go/internal/service"

type Application struct {
	UserService *service.UserService        
}
