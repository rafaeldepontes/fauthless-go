package errorhandler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

var ErrorUserNotFound = errors.New("User not found")
var ErrorInvalidMethod = errors.New("Invalid method")
var ErrorUserAlreadyExists = errors.New("User already exist")
var ErrorInvalidUsernameOrPassword = errors.New("Invalid username or password")
var ErrorInvalidTokenSigningMethod = errors.New("Invalid token signing method")
var ErrorParsingToken = errors.New("Error parsing token")
var ErrorInvalidTokenClaim = errors.New("Invalid token claims")
var ErrorCreatingToken = errors.New("Error while creating token")

var ErrorIdIsRequired = errors.New("Identifier is required")
var ErrorUsernameIsRequired = errors.New("Username is required")
var ErrorPasswordIsRequired = errors.New("Password is required")
var ErrorAgeIsRequired = errors.New("Age is required")

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
	UnauthroizedErrorHandler = func(w http.ResponseWriter) {
		writeError(w, "Unauthroized", http.StatusUnauthorized, "")
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
