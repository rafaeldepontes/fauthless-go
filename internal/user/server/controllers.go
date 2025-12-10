package server

import (
	"net/http"

	"github.com/rafaeldepontes/fauthless-go/internal/user"
)

type userController struct {
	service *user.Service
}

func NewUserController(s *user.Service) user.Controller {
	return &userController{
		service: s,
	}
}

func (c *userController) ListAllHashedCursor(w http.ResponseWriter, r *http.Request) {
	(*c.service).FindAllUsersHashedCursorPagination(w, r)
}

func (c *userController) ListAllCursor(w http.ResponseWriter, r *http.Request) {
	(*c.service).FindAllUsersCursorPagination(w, r)
}

func (c *userController) ListAllOffset(w http.ResponseWriter, r *http.Request) {
	(*c.service).FindAllUsersOffSetPagination(w, r)
}

func (c *userController) FindById(w http.ResponseWriter, r *http.Request) {
	(*c.service).FindUserById(w, r)
}

func (c *userController) UpdateDetails(w http.ResponseWriter, r *http.Request) {
	(*c.service).UpdateUserDetails(w, r)
}

func (c *userController) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	(*c.service).DeleteAccount(w, r)
}
