package errorhandler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	ErrInvalidUsernameOrPassword = errors.New("Invalid username or password")
	ErrInvalidTokenSigningMethod = errors.New("Invalid token signing method")
	ErrInvalidTokenSignature     = errors.New("Token signature is invalid")
	ErrInvalidExpiredToken       = errors.New("Invalid token, token already expired")
	ErrUsernameIsRequired        = errors.New("Username is required")
	ErrPasswordIsRequired        = errors.New("Password is required")
	ErrInvalidTokenClaim         = errors.New("Invalid token claims")
	ErrUserAlreadyExists         = errors.New("User already exist")
	ErrTokenNotValidYet          = errors.New("Token is not valid yet")
	ErrInvalidCSRFToken          = errors.New("CSRF token missing")
	ErrMalformedToken            = errors.New("Token is malformed")
	ErrInvalidSession            = errors.New("Session not found")
	ErrAgeIsRequired             = errors.New("Age is required")
	ErrInvalidMethod             = errors.New("Invalid method")
	ErrCreatingToken             = errors.New("Error while creating token")
	ErrEqualUsername             = errors.New("The new username should be different from the actual")
	ErrIdIsRequired              = errors.New("Identifier is required")
	ErrUserNotFound              = errors.New("User not found")
	ErrParsingToken              = errors.New("Error parsing token")
	ErrInvalidToken              = errors.New("Token missing or invalid")
	ErrTokenRevoked              = errors.New("Session token revoked")
	ErrInvalidType               = errors.New("Unsupported type")
	ErrInvalidId                 = errors.New("Invalid username, needs to be your own")
	ErrEqualAge                  = errors.New("The new age should be different from the actual")
)

const BrazilianDateTimeFormat = "02/01/2006 15:04:05"

type Error struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	Path      string `json:"path,omitempty"`
	Status    int    `json:"status"`
}

var (
	BadRequestErrorHandler = func(w http.ResponseWriter, err error, path string) {
		writeError(w, err.Error(), http.StatusBadRequest, path)
	}
	InternalErrorHandler = func(w http.ResponseWriter) {
		writeError(w, "An unexpected Error Occurred.", http.StatusInternalServerError, "")
	}
	UnauthroizedErrorHandler = func(w http.ResponseWriter, err error) {
		writeError(w, err.Error(), http.StatusUnauthorized, "")
	}
	ForbiddenErrorHandler = func(w http.ResponseWriter, err error) {
		writeError(w, err.Error(), http.StatusForbidden, "")
	}
	RequestErrorHandler = func(w http.ResponseWriter, err error, status int, path string) {
		writeError(w, err.Error(), status, path)
	}
)

func writeError(w http.ResponseWriter, message string, status int, path string) {
	var timestamp string = time.Now().Format(BrazilianDateTimeFormat)
	resp := Error{
		Status:    status,
		Message:   message,
		Path:      path,
		Timestamp: timestamp,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Error(err)
	}
}
