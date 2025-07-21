package server

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Server struct {
	db *pgxpool.Pool
}

// Initialize server and database
func NewServer() (*Server, error) {
	dbpool, err := pgxpool.Connect(context.Background(), "postgres://user:password@localhost:5432/docs")
	if err != nil {
		return nil, err
	}

	return &Server{db: dbpool}, nil
}
