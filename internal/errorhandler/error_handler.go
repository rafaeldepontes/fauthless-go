package errorhandler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

var ErrorUsernameNotFound = errors.New("User not found")
var ErrorInvalidMethod = errors.New("Invalid method")
var ErrorUserAlreadyExists = errors.New("User already exist")

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
	RequestErrorHandler = func(w http.ResponseWriter, err error, status int, path string) {
		writeError(w, err.Error(), status, path)
	}
	InternalErrorHandler = func(w http.ResponseWriter) {
		writeError(w, "An unexpected Error Occurred.", http.StatusInternalServerError, "")
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
