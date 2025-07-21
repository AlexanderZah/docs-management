package docs

import (
	"context"
)

type Repository interface {
	SaveDocument(ctx context.Context, doc *Document) error
}

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) UploadDocument(ctx context.Context, doc *Document) error {
	// Можно добавить валидацию или логику
	return s.repo.SaveDocument(ctx, doc)
}
