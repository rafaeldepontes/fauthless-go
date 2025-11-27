package service

import (
	"net/http"

	"github.com/rafaeldepontes/auth-go/internal/repository"
)

type UserService struct {
	userRepository *repository.UserRepository
}

// NewUserService initialize a new UserService containing a UserRepository
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepo,
	}
}

// FindAllUsers list all the users without a filter and returns each
// one with pagination and a few datas missing for LGPD
func (us *UserService) FindAllUsers(w http.ResponseWriter, r *http.Request) {
	// TODO: WIP
}
