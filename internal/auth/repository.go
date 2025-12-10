package auth

import "github.com/rafaeldepontes/fauthless-go/internal/domain"

type Repository interface {
	CreateSession(session *domain.Session) (string, error)
	FindSessionById(id string) (*domain.Session, error)
	RevokeSession(id string) error
	DeleteSession(id string) error
}
