package service

import (
	"net/http"

	"github.com/rafaeldepontes/auth-go/internal/repository"
	log "github.com/sirupsen/logrus"
)

type AuthService struct {
	userRepository *repository.UserRepository
	Logger         *log.Logger
}

// NewAuthService initialize a new AuthService containing a UserRepository for
// login and register operations ONLY.
func NewAuthService(userRepo *repository.UserRepository, logg *log.Logger) *AuthService {
	return &AuthService{
		userRepository: userRepo,
		Logger:         logg,
	}
}

func (as AuthService) Register(w http.ResponseWriter, r *http.Request) {

}

func (as AuthService) Login(w http.ResponseWriter, r *http.Request) {

}
