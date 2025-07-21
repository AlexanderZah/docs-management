package user

import (
	"errors"
	"regexp"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	CreateUser(u *User) error
	Exists(login string) (bool, error)
}

type Service struct {
	repo        Repository
	adminSecret string
}

func NewService(r Repository, adminToken string) *Service {
	return &Service{repo: r, adminSecret: adminToken}
}

func (s *Service) Register(adminToken, login, password string) (*User, error) {
	if adminToken != s.adminSecret {
		return nil, errors.New("invalid admin token")
	}

	if err := validateLogin(login); err != nil {
		return nil, err
	}
	if err := validatePassword(password); err != nil {
		return nil, err
	}

	exists, err := s.repo.Exists(login)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("user already exists")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:       uuid.NewString(),
		Login:    login,
		Password: string(hashed),
	}

	return user, s.repo.CreateUser(user)
}

func validateLogin(login string) error {
	if len(login) < 8 || !regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(login) {
		return errors.New("login must be at least 8 characters and alphanumeric")
	}
	return nil
}

func validatePassword(p string) error {
	if len(p) < 8 {
		return errors.New("password too short")
	}
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(p)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(p)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(p)
	hasSymbol := regexp.MustCompile(`[^a-zA-Z0-9]`).MatchString(p)

	if !hasLower || !hasUpper || !hasDigit || !hasSymbol {
		return errors.New("password must include upper, lower, digit, and symbol")
	}
	return nil
}
