package service

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/rafaeldepontes/auth-go/internal/domain"
	"github.com/rafaeldepontes/auth-go/internal/errorhandler"
	"github.com/rafaeldepontes/auth-go/internal/repository"
	"github.com/rafaeldepontes/auth-go/internal/token"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const (
	Cost         = 16
	Token_Length = 32
)

type AuthService struct {
	jwtMaker       *token.JwtBuilder
	userRepository *repository.UserRepository
	Logger         *log.Logger
}

// NewAuthService initialize a new AuthService containing a UserRepository for
// login and register operations ONLY.
func NewAuthService(userRepo *repository.UserRepository, logg *log.Logger, secretKey string) *AuthService {
	return &AuthService{
		userRepository: userRepo,
		Logger:         logg,
		jwtMaker:       token.NewJwtBuilder(secretKey),
	}
}

// Register is a generic register system that can be use in any case,
// it doesnt returns nothing and only insert a new user into the database
// after a bunch of validations.
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	as.Logger.Infoln("The user registered successfully.")
}

func (as AuthService) LoginCookieBased(w http.ResponseWriter, r *http.Request) {
	var user *repository.User = loginFlow(&as, w, r)
	token := token.CookieBased{}

	sessionToken := token.GenerateToken(Token_Length)
	csrfToken := token.GenerateToken(Token_Length)

	err := as.userRepository.SetUserToken(sessionToken, csrfToken, *user.Id)
	if err != nil {
		as.Logger.Errorf("An error occurred: %v", err)
		errorhandler.InternalErrorHandler(w)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(30 * time.Minute),
		HttpOnly: true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:    "crsf_token",
		Value:   csrfToken,
		Expires: time.Now().Add(30 * time.Minute),
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	as.Logger.Infoln("The user logged in successfully.")
}

func (as AuthService) LoginJwtBased(w http.ResponseWriter, r *http.Request) {
	var user *repository.User = loginFlow(&as, w, r)
	var maker *token.JwtBuilder = as.jwtMaker

	id, username := user.Id, user.Username
	as.Logger.Infoln(id, username)

	token, _, err := maker.GenerateToken(*id, *username, 15*time.Minute)
	if err != nil {
		as.Logger.Errorf("An error occurred: %v", err)
		errorhandler.InternalErrorHandler(w)
		return
	}

	userResponse := domain.TokenResponse{
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userResponse)
}

func (as AuthService) LoginJwtRefreshBased(w http.ResponseWriter, r *http.Request) {

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

func loginFlow(as *AuthService, w http.ResponseWriter, r *http.Request) *repository.User {
	as.Logger.Infoln("Trying to login user")

	if r.Method != http.MethodPost {
		as.Logger.Errorf("An error occurred: %v", errorhandler.ErrorInvalidMethod)
		errorhandler.BadRequestErrorHandler(w, errorhandler.ErrorInvalidMethod, r.URL.Path)
		return nil
	}

	var user domain.UserLogin
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		as.Logger.Errorf("An error occurred: %v", err)
		errorhandler.InternalErrorHandler(w)
		return nil
	}

	var userInTheDatabase *repository.User
	userInTheDatabase, err = as.userRepository.FindUserByUsername(user.Username)
	if err != nil {
		as.Logger.Errorf("An error occurred: %v", errorhandler.ErrorUserNotFound)
		errorhandler.BadRequestErrorHandler(w, errorhandler.ErrorUserNotFound, r.URL.Path)
		return nil
	}

	password := *userInTheDatabase.HashedPassword
	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password))
	if err != nil {
		as.Logger.Errorf("An error occurred: %v", errorhandler.ErrorInvalidUsernameOrPassword)
		errorhandler.BadRequestErrorHandler(w, errorhandler.ErrorInvalidUsernameOrPassword, r.URL.Path)
		return nil
	}

	as.Logger.Infoln("Valid user, following the next steps...")
	return userInTheDatabase
}
