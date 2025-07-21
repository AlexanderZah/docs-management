package repository

import (
	"context"
	"errors"

	"github.com/AlexanderZah/docs-management/internal/user"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{db: pool}
}

func (r *Repo) Exists(login string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(context.Background(), `
        SELECT EXISTS (SELECT 1 FROM users WHERE login = $1)
    `, login).Scan(&exists)
	return exists, err
}

func (r *Repo) CreateUser(u *user.User) error {

	_, err := r.db.Exec(context.Background(), `
        INSERT INTO users (id, login, password) 
        VALUES ($1, $2, $3)
    `, u.ID, u.Login, u.Password)

	if err != nil {
		if pgErr, ok := err.(*pgx.PgError); ok && pgErr.Code == "23505" {
			return errors.New("user already exists")
		}
	}
	return err
}

func (r *Repo) FindByLogin(login string) (*user.User, error) {
	var u user.User
	err := r.db.QueryRow(context.Background(), `
        SELECT id, login, password FROM users WHERE login = $1
    `, login).Scan(&u.ID, &u.Login, &u.Password)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
