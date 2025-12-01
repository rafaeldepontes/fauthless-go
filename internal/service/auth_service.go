package service

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/rafaeldepontes/auth-go/internal/database/repository"
	"github.com/rafaeldepontes/auth-go/internal/domain"
	"github.com/rafaeldepontes/auth-go/internal/errorhandler"
	"github.com/rafaeldepontes/auth-go/internal/storage"
	"github.com/rafaeldepontes/auth-go/internal/token"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const (
	Cost         = 16
	Token_Length = 32
)

type AuthService struct {
	jwtMaker          *token.JwtBuilder
	userRepository    *repository.UserRepository
	sessionRepository *repository.SessionRepository
	Logger            *log.Logger
	Cache             *storage.Caches
}

// NewAuthService initialize a new AuthService containing a UserRepository for
// login and register operations ONLY.
func NewAuthService(userRepo *repository.UserRepository, sessionRepo *repository.SessionRepository, logg *log.Logger, secretKey string, cache *storage.Caches) *AuthService {
	return &AuthService{
		userRepository:    userRepo,
		sessionRepository: sessionRepo,
		Logger:            logg,
		jwtMaker:          token.NewJwtBuilder(secretKey),
		Cache:             cache,
	}
}

// Register is a generic register system that can be use in any case,
// it doesnt returns nothing and only insert a new user into the database
// after a bunch of validations.
func (s *AuthService) Register(w http.ResponseWriter, r *http.Request) {
	s.Logger.Infoln("Registering a new user")

	if r.Method != http.MethodPost {
		s.Logger.Errorf("An error occurred: %v", errorhandler.ErrInvalidMethod)
		errorhandler.BadRequestErrorHandler(w, errorhandler.ErrInvalidMethod, r.URL.Path)
		return
	}

	var user repository.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		s.Logger.Errorf("An error occurred: %v", err)
		errorhandler.InternalErrorHandler(w)
		return
	}

	if ok, err := isValidUser(&user, s); !ok {
		s.Logger.Errorf("An error occurred: %v", err)
		errorhandler.BadRequestErrorHandler(w, err, r.URL.Path)
		return
	}

	password := user.HashedPassword

	var hashedPassword []byte
	hashedPassword, err = bcrypt.GenerateFromPassword([]byte(*password), Cost)
	if err != nil {
		s.Logger.Errorf("An error occurred: %v", err)
		errorhandler.InternalErrorHandler(w)
		return
	}

	*password = string(hashedPassword)

	err = s.userRepository.RegisterUser(&user)
	if err != nil {
		s.Logger.Errorf("An error occurred: %v", err)
		errorhandler.InternalErrorHandler(w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	s.Logger.Infoln("The user registered successfully.")
}

// LoginCookieBased uses the cookie authorization flow, creating a token that needs to be
// in the request cookie.
func (s *AuthService) LoginCookieBased(w http.ResponseWriter, r *http.Request) {
	var user *repository.User = loginFlow(s, w, r)
	if user == nil {
		return
	}

	token := token.CookieBased{}

	sessionToken := token.GenerateToken(Token_Length)
	csrfToken := token.GenerateToken(Token_Length)

	err := s.userRepository.SetUserToken(sessionToken, csrfToken, *user.Id)
	if err != nil {
		s.Logger.Errorf("An error occurred: %v", err)
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
		Name:    "csrf_token",
		Value:   csrfToken,
		Expires: time.Now().Add(30 * time.Minute),
	})

	durationInt, _ := strconv.Atoi(os.Getenv("TOKEN_DURATION"))
	var duration time.Duration = time.Duration(durationInt)

	userCache := s.Cache.UserCache
	userCache.Set(sessionToken, *user.Username, time.Now().Add(duration*time.Minute))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	s.Logger.Infoln("The user logged in successfully.")
}

// LoginJwtBased uses the Jwt method to create a access token, with it
// all the features are available until it expires.
func (s *AuthService) LoginJwtBased(w http.ResponseWriter, r *http.Request) {
	var user *repository.User = loginFlow(s, w, r)
	if user == nil {
		return
	}

	var maker *token.JwtBuilder = s.jwtMaker

	durationInt, _ := strconv.Atoi(os.Getenv("TOKEN_DURATION"))
	var duration time.Duration = time.Duration(durationInt)
	token, _, err := maker.GenerateToken(*user.Id, *user.Username, duration*time.Minute)
	if err != nil {
		s.Logger.Errorf("An error occurred: %v", err)
		errorhandler.InternalErrorHandler(w)
		return
	}

	tokenCache := s.Cache.TokenCache
	invalid := false
	tokenCache.Set(token, invalid, time.Now().Add(duration*time.Minute))

	userResponse := domain.TokenResponse{
		Token: token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userResponse)
}

// LoginJwtRefreshBased uses the Jwt method to create a access token, with it
// all the features are available until it expires, but it cames with a refresh
// token that can be used in another call to gain access again until the refresh
// one expires...
func (s *AuthService) LoginJwtRefreshBased(w http.ResponseWriter, r *http.Request) {
	var user *repository.User = loginFlow(s, w, r)
	if user == nil {
		return
	}

	var maker *token.JwtBuilder = s.jwtMaker
	accessToken, accessClaims, err := generateAccessToken(maker, *user.Id, *user.Username, time.Minute)
	if err != nil {
		s.Logger.Errorf("An error occurred: %v", err)
		errorhandler.InternalErrorHandler(w)
		return
	}

	refreshToken, refreshClaims, err := generateTokenRefresh(maker, *user.Id, *user.Username, time.Hour)
	if err != nil {
		s.Logger.Errorf("An error occurred: %v", err)
		errorhandler.InternalErrorHandler(w)
		return
	}

	sessionId, err := s.sessionRepository.CreateSession(&repository.Session{
		Id:           refreshClaims.RegisteredClaims.ID,
		Username:     *user.Username,
		RefreshToken: refreshToken,
		IsRevoked:    false,
		ExpiresAt:    refreshClaims.ExpiresAt.Time,
	})
	if err != nil {
		s.Logger.Errorf("An error occurred: %v", err)
		errorhandler.InternalErrorHandler(w)
		return
	}

	trResponse := domain.TokenRefreshResponse{
		SessionId:             sessionId,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresAt:  accessClaims.ExpiresAt.Time,
		RefreshTokenExpiresAt: refreshClaims.ExpiresAt.Time,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(trResponse)
}

// RenewAccessToken accepts a json body for the request, in it should have the
// refresh token available at the login call. After called and with a proper token
// it gives another access token for futher uses.
func (s *AuthService) RenewAccessToken(w http.ResponseWriter, r *http.Request) {
	var maker *token.JwtBuilder = s.jwtMaker

	var req domain.RenewAccessTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.Logger.Errorf("An error occurred: %v", err)
		errorhandler.InternalErrorHandler(w)
		return
	}

	refreshClaims, err := s.jwtMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		s.Logger.Errorf("An error occurred: %v", err)
		errorhandler.BadRequestErrorHandler(w, err, r.URL.Path)
		return
	}

	var session *repository.Session
	session, err = s.sessionRepository.FindSessionById(refreshClaims.ID)
	if err != nil {
		s.Logger.Errorf("An error occurred: %v", err)
		errorhandler.BadRequestErrorHandler(w, errorhandler.ErrInvalidToken, r.URL.Path)
		return
	}

	if session.IsRevoked {
		s.Logger.Errorf("An error occurred: %v", errorhandler.ErrTokenRevoked)
		errorhandler.BadRequestErrorHandler(w, errorhandler.ErrTokenRevoked, r.URL.Path)
		return
	}

	if session.Username != refreshClaims.Username {
		s.Logger.Errorf("An error occurred: %v", errorhandler.ErrTokenRevoked)
		errorhandler.ForbiddenErrorHandler(w, errorhandler.ErrTokenRevoked)
		return
	}

	accessToken, accessClaims, err := generateAccessToken(maker, refreshClaims.Id, refreshClaims.Username, time.Minute)
	if err != nil {
		s.Logger.Errorf("An error occurred: %v", err)
		errorhandler.InternalErrorHandler(w)
		return
	}

	tkResponse := &domain.RenewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessClaims.ExpiresAt.Time,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tkResponse)
}

// RevokeSession disable a refresh token, preventing futher requests.
func (s *AuthService) RevokeSession(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		s.Logger.Errorf("An error occurred: %v", errorhandler.ErrIdIsRequired)
		errorhandler.BadRequestErrorHandler(w, errorhandler.ErrIdIsRequired, r.URL.Path)
		return
	}

	err := s.sessionRepository.RevokeSession(id)
	if err != nil {
		s.Logger.Errorf("An error occurred: %v", err)
		errorhandler.BadRequestErrorHandler(w, errorhandler.ErrInvalidSession, r.URL.Path)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func isValidUser(newUser *repository.User, s *AuthService) (bool, error) {
	if username := newUser.Username; *username == "" {
		return false, errorhandler.ErrUsernameIsRequired
	}

	if password := newUser.HashedPassword; *password == "" {
		return false, errorhandler.ErrPasswordIsRequired
	}

	if age := newUser.Age; *age == 0 {
		return false, errorhandler.ErrAgeIsRequired
	}

	user, err := s.userRepository.FindUserByUsername(*newUser.Username)
	if user != nil {
		s.Logger.Errorf("An error occurred: %v\n", err)
		return false, errorhandler.ErrUserAlreadyExists
	}

	return true, nil
}

func loginFlow(s *AuthService, w http.ResponseWriter, r *http.Request) *repository.User {
	s.Logger.Infoln("Trying to login user")

	if r.Method != http.MethodPost {
		s.Logger.Errorf("An error occurred: %v", errorhandler.ErrInvalidMethod)
		errorhandler.BadRequestErrorHandler(w, errorhandler.ErrInvalidMethod, r.URL.Path)
		return nil
	}

	var user domain.UserLogin
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		s.Logger.Errorf("An error occurred: %v", err)
		errorhandler.InternalErrorHandler(w)
		return nil
	}

	var userInTheDatabase *repository.User
	userInTheDatabase, err = s.userRepository.FindUserByUsername(user.Username)
	if err != nil {
		s.Logger.Errorf("An error occurred: %v", errorhandler.ErrUserNotFound)
		errorhandler.BadRequestErrorHandler(w, errorhandler.ErrUserNotFound, r.URL.Path)
		return nil
	}

	password := *userInTheDatabase.HashedPassword
	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password))
	if err != nil {
		s.Logger.Errorf("An error occurred: %v", errorhandler.ErrInvalidUsernameOrPassword)
		errorhandler.BadRequestErrorHandler(w, errorhandler.ErrInvalidUsernameOrPassword, r.URL.Path)
		return nil
	}

	s.Logger.Infoln("Valid user, following the next steps...")
	return userInTheDatabase
}

func generateTokenRefresh(maker *token.JwtBuilder, id uint, username string, timer time.Duration) (string, *token.UserClaims, error) {
	return generateToken(maker, id, username, 24, timer)
}

func generateAccessToken(maker *token.JwtBuilder, id uint, username string, timer time.Duration) (string, *token.UserClaims, error) {
	durationInt, _ := strconv.Atoi(os.Getenv("TOKEN_DURATION"))
	var duration time.Duration = time.Duration(durationInt)
	return generateToken(maker, id, username, duration, timer)
}

func generateToken(maker *token.JwtBuilder, id uint, username string, timer time.Duration, duration time.Duration) (string, *token.UserClaims, error) {
	token, userClaims, err := maker.GenerateToken(id, username, duration*timer)
	if err != nil {
		return "", nil, err
	}

	return token, userClaims, nil
}
