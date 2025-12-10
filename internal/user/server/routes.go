package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/rafaeldepontes/fauthless-go/internal/user"
)

func MapUserRoutes(route *chi.Router, controller *user.Controller) {
	// I could have done this in the same request, but for the learning purposes,
	// I'm doing it separately.
	(*route).Get("/users/hashed-cursor-pagination", (*controller).ListAllHashedCursor)
	(*route).Get("/users/cursor-pagination", (*controller).ListAllCursor)
	(*route).Get("/users/offset-pagination", (*controller).ListAllOffset)

	(*route).Get("/users/{id}", (*controller).FindById)
}

func MapUserRoutesJwt(route *chi.Router, controller *user.Controller) {
	(*route).Patch("/users/{username}", (*controller).UpdateDetails)
	(*route).Delete("/users/{username}", (*controller).DeleteAccount)
}
