package service

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/rafaeldepontes/auth-go/internal/errorhandler"
	"github.com/rafaeldepontes/auth-go/internal/pagination"
	"github.com/rafaeldepontes/auth-go/internal/repository"
	log "github.com/sirupsen/logrus"
)

type UserService struct {
	userRepository *repository.UserRepository
	Logger         *log.Logger
}

// NewUserService initialize a new UserService containing a UserRepository.
func NewUserService(userRepo *repository.UserRepository, logg *log.Logger) *UserService {
	return &UserService{
		userRepository: userRepo,
		Logger:         logg,
	}
}

// FindAllUsers list all the users without a filter and returns each
// one with pagination and a few datas missing for LGPD.
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

	var users []repository.User
	users, totalRecords, err := us.userRepository.FindAllUsers(size, currentPage)
	if err != nil {
		errorhandler.BadRequestErrorHandler(w, err, r.URL.Path)
		us.Logger.Errorf("An error occurred: %v", err)
	}

	pageModel := pagination.NewPagination(users, uint(currentPage), uint(totalRecords), uint(size))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	us.Logger.Infof("Found %v users from a total of %v", len(users), totalRecords)

	json.NewEncoder(w).Encode(pageModel)
}

// FindUserById list an user by his id and returns a none
// pagination result and a few datas missing for LGPD.
func (us UserService) FindUserById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	us.Logger.Infof("Listing user by id - %v", idStr)

	if idStr == "" {
		errorhandler.BadRequestErrorHandler(w, errorhandler.ErrorIdIsRequired, r.URL.Path)
		us.Logger.Errorf("An error occurred: %v", errorhandler.ErrorIdIsRequired)
	}

	pathId, _ := strconv.Atoi(idStr)
	id := uint(pathId)

	var user repository.User
	user, err := us.userRepository.FindUserById(id)
	if err != nil {
		errorhandler.BadRequestErrorHandler(w, errorhandler.ErrorUsernameNotFound, r.URL.Path)
		us.Logger.Errorf("An error occurred: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(user)
}

func (us UserService) UpdateUserDetails(w http.ResponseWriter, r *http.Request) {

}

func (us UserService) DeleteAccount(w http.ResponseWriter, r *http.Request) {

}
