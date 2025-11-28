package service

import (
	"encoding/json"
	"net/http"

	"github.com/rafaeldepontes/auth-go/internal/domain"
	"github.com/rafaeldepontes/auth-go/internal/errorhandler"
	"github.com/rafaeldepontes/auth-go/internal/repository"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const Cost = 16

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

func (as *AuthService) Register(w http.ResponseWriter, r *http.Request) {
	as.Logger.Infoln("Registering a new user")

	if r.Method != http.MethodPost {
		as.Logger.Errorf("An error occurred: %v", errorhandler.ErrorInvalidMethod)
		errorhandler.BadRequestErrorHandler(w, errorhandler.ErrorInvalidMethod, r.URL.Path)
		return
	}

	var user repository.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		as.Logger.Errorf("An error occurred: %v", err)
		errorhandler.InternalErrorHandler(w)
		return
	}

	if ok, err := isValidUser(&user, as); !ok {
		as.Logger.Errorf("An error occurred: %v", err)
		errorhandler.BadRequestErrorHandler(w, err, r.URL.Path)
		return
	}

	password := user.HashedPassword

	var hashedPassword []byte
	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(*password), Cost)
	if err != nil {
		as.Logger.Errorf("An error occurred: %v", err)
		errorhandler.InternalErrorHandler(w)
		return
	}

	*password = string(hashedPassword)

	err = as.userRepository.RegisterUser(&user)
	if err != nil {
		as.Logger.Errorf("An error occurred: %v", err)
		errorhandler.InternalErrorHandler(w)
		return
	}

	as.Logger.Infoln("The user registered successfully.")
}

func (as AuthService) Login(w http.ResponseWriter, r *http.Request) {
	as.Logger.Infoln("Trying to login user")

	if r.Method != http.MethodPost {
		as.Logger.Errorf("An error occurred: %v", errorhandler.ErrorInvalidMethod)
		errorhandler.BadRequestErrorHandler(w, errorhandler.ErrorInvalidMethod, r.URL.Path)
		return
	}

	var user domain.UserLogin
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		as.Logger.Errorf("An error occurred: %v", err)
		errorhandler.InternalErrorHandler(w)
		return
	}

	var userInTheDatabase *repository.User
	userInTheDatabase, err = as.userRepository.FindUserByUsername(user.Username)
	if err != nil {
		as.Logger.Errorf("An error occurred: %v", errorhandler.ErrorUserNotFound)
		errorhandler.BadRequestErrorHandler(w, errorhandler.ErrorUserNotFound, r.URL.Path)
		return
	}

	password := *userInTheDatabase.HashedPassword
	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password))
	if err != nil {
		as.Logger.Errorf("An error occurred: %v", errorhandler.ErrorInvalidUsernameOrPassword)
		errorhandler.BadRequestErrorHandler(w, errorhandler.ErrorInvalidUsernameOrPassword, r.URL.Path)
		return
	}

	as.Logger.Infoln("The user logged in successfully.")
}

func isValidUser(newUser *repository.User, as *AuthService) (bool, error) {
	if username := newUser.Username; *username == "" {
		return false, errorhandler.ErrorUsernameIsRequired
	}

	if password := newUser.HashedPassword; *password == "" {
		return false, errorhandler.ErrorPasswordIsRequired
	}

	if age := newUser.Age; *age == 0 {
		return false, errorhandler.ErrorAgeIsRequired
	}

	emptyUser := repository.User{}
	user, _ := as.userRepository.FindUserByUsername(*newUser.Username)
	if *user != emptyUser {
		return false, errorhandler.ErrorUserAlreadyExists
	}

	return true, nil
}
