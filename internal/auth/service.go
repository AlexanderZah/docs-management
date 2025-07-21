package auth

import (
	"errors"

	"github.com/AlexanderZah/docs-management/internal/user"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type TokenStore interface {
	Save(login, token string) error
	GetLogin(token string) (string, error)
	Delete(token string) error
}

type Service struct {
	users user.Repository
	store TokenStore
}

func NewService(users user.Repository, store TokenStore) *Service {
	return &Service{users: users, store: store}
}

func (s *Service) Authenticate(login, password string) (string, error) {
	u, err := s.users.FindByLogin(login)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	token := uuid.NewString()
	err = s.store.Save(login, token)
	return token, err
}
