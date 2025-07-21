package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/AlexanderZah/docs-management/internal/user"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Repo struct {
	db *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{db: pool}
}

func (r *Repo) Exists(ctx context.Context, login string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `
        SELECT EXISTS (SELECT 1 FROM users WHERE login = $1)
    `, login).Scan(&exists)
	return exists, err
}

func (r *Repo) CreateUser(ctx context.Context, u *user.User) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback(ctx)
	const query = `
        INSERT INTO users (login, password) 
        VALUES ($1, $2)
    `
	_, err = tx.Exec(ctx, query, u.Login, u.Password)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return errors.New("user already exists")
		}
	}
	return tx.Commit(ctx)
}

func (r *Repo) FindByLogin(ctx context.Context, login string) (*user.User, error) {
	var u user.User
	err := r.db.QueryRow(ctx, `
        SELECT id, login, password FROM users WHERE login = $1
    `, login).Scan(&u.ID, &u.Login, &u.Password)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
