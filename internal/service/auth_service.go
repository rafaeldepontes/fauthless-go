package service

import (
	"net/http"

	"github.com/rafaeldepontes/auth-go/internal/errorhandler"
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
	username := r.FormValue("username")
	if username == "" {
		as.Logger.Errorf("An error occurred: %v", errorhandler.ErrorUsernameIsRequired)
		errorhandler.BadRequestErrorHandler(w, errorhandler.ErrorUsernameIsRequired, r.URL.Path)
		return
	}

	// var user = repository.User{
	// 	Username: ,
	// }

	// as.userRepository.RegisterUser(&user)
}

func (as AuthService) Login(w http.ResponseWriter, r *http.Request) {

}
