package service

import (
	"encoding/json"
	"github.com/rafaeldepontes/fauthless-go/internal/domain"
	"github.com/rafaeldepontes/fauthless-go/internal/errorhandler"
	"github.com/rafaeldepontes/fauthless-go/internal/token"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strconv"
	"time"
)

func loginFlow(s *authService, w http.ResponseWriter, r *http.Request) *domain.User {
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

	var userInTheDatabase *domain.User
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

func generateTokenRefresh(maker *token.JwtBuilder, id int64, username string, timer time.Duration) (string, *token.UserClaims, error) {
	return generateToken(maker, id, username, 24, timer)
}

func generateAccessToken(maker *token.JwtBuilder, id int64, username string, timer time.Duration) (string, *token.UserClaims, error) {
	durationInt, _ := strconv.Atoi(os.Getenv("TOKEN_DURATION"))
	var duration time.Duration = time.Duration(durationInt)
	return generateToken(maker, id, username, duration, timer)
}

func generateToken(maker *token.JwtBuilder, id int64, username string, timer time.Duration, duration time.Duration) (string, *token.UserClaims, error) {
	token, userClaims, err := maker.GenerateToken(id, username, duration*timer)
	if err != nil {
		return "", nil, err
	}

	return token, userClaims, nil
}
