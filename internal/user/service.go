package user

import (
	"context"
	"errors"
	"regexp"

	er "github.com/AlexanderZah/docs-management/internal/myerrors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	CreateUser(ctx context.Context, u *User) error
	Exists(ctx context.Context, login string) (bool, error)
	FindByLogin(ctx context.Context, login string) (*User, error)
	GetSession(ctx context.Context, login string) (string, error)
	SaveSession(ctx context.Context, token string, login string) error
	DeleteSession(ctx context.Context, token string) error
	GetSessionByToken(ctx context.Context, token string) (string, error)
}

type Service struct {
	repo        Repository
	adminSecret string
}

func NewService(r Repository, adminToken string) *Service {
	return &Service{repo: r, adminSecret: adminToken}
}

func (s *Service) Register(ctx context.Context, adminToken, login, password string) (*User, error) {
	if adminToken != s.adminSecret {
		return nil, errors.New("invalid admin token")
	}

	if err := validateLogin(login); err != nil {
		return nil, err
	}
	if err := validatePassword(password); err != nil {
		return nil, err
	}

	exists, err := s.repo.Exists(ctx, login)
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
		Login:    login,
		Password: string(hashed),
	}

	return user, s.repo.CreateUser(ctx, user)
}

func (s *Service) Login(ctx context.Context, login, password string) (string, error) {
	u, err := s.repo.FindByLogin(ctx, login)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}
	token, err := s.repo.GetSession(ctx, login)
	if err != nil && !errors.Is(err, er.ErrSessionNotFound) {
		return "", errors.New(err.Error())
	}
	if errors.Is(err, er.ErrSessionNotFound) {
		token = uuid.NewString()
		s.repo.SaveSession(ctx, token, login)
	}

	return token, nil
}

func (s *Service) Logout(ctx context.Context, token string) error {
	return s.repo.DeleteSession(ctx, token)
}

func (s *Service) GetSessionByToken(ctx context.Context, token string) (string, error) {
	return s.repo.GetSessionByToken(ctx, token)
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
