package user

import "github.com/rafaeldepontes/fauthless-go/internal/domain"

type Repository interface {
	FindAllUsersCursor(cursor int64, size int) ([]domain.User, int64, error)
	FindAllUsers(size, page int) ([]domain.User, int, error)
	FindUserById(id int64) (*domain.User, error)
	FindUserByUsername(username string) (*domain.User, error)
	RegisterUser(u *domain.User) error
	SetUserToken(token, csrfToken string, userId int64) error
	UpdateUserDetails(user *domain.User) error
	DeleteAccount(username string) error
}
