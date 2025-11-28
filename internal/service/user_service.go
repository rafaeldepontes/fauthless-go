package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/rafaeldepontes/auth-go/internal/errorhandler"
	"github.com/rafaeldepontes/auth-go/internal/repository"
	log "github.com/sirupsen/logrus"
)

type UserService struct {
	Logger         *log.Logger
	userRepository *repository.UserRepository
}

// NewUserService initialize a new UserService containing a UserRepository
func NewUserService(userRepo *repository.UserRepository, logg *log.Logger) *UserService {
	return &UserService{
		Logger:         logg,
		userRepository: userRepo,
	}
}

// FindAllUsers list all the users without a filter and returns each
// one with pagination and a few datas missing for LGPD
func (us *UserService) FindAllUsers(w http.ResponseWriter, r *http.Request) {
	us.Logger.Infoln("Listing all the users in the database...")

	sizeStr := r.URL.Query().Get("size")
	if sizeStr == "" {
		sizeStr = "25"
	}
	size, _ := strconv.Atoi(sizeStr)

	currentPageStr := r.URL.Query().Get("page")
	if currentPageStr == "" {
		currentPageStr = "1"
	}
	currentPage, _ := strconv.Atoi(currentPageStr)
	currentPage-- // If page is the first one, in the URL should be "1" and in the code should be 0 for the offset...

	var users []repository.User
	users, err := us.userRepository.FindAllUsers(size, currentPage)
	if err != nil {
		errorhandler.BadRequestErrorHandler(w, err, r.URL.Path)
		log.Errorf("An error occurred: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	us.Logger.Infof("Found %v users", len(users))

	//TODO: IMPLEMENT PAGINATION FOR THIS ENDPOINT
	json.NewEncoder(w).Encode(users)
}
