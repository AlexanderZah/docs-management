package docs

import (
	"context"
)

type Repository interface {
	SaveDocument(ctx context.Context, doc *Document) error
	GetDocuments(ctx context.Context, token, login, key, value string, limit int) ([]Document, error)
	GetDocumentByID(ctx context.Context, id int32, token string) (*Document, error)
	DeleteDocument(ctx context.Context, id int32, token string) error
}

type Service struct {
	repo Repository
}

func NewService(r Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) UploadDocument(ctx context.Context, doc *Document) error {
	return s.repo.SaveDocument(ctx, doc)
}

func (s *Service) GetDocuments(ctx context.Context, token, login, key, value string, limit int) ([]Document, error) {
	return s.repo.GetDocuments(ctx, token, login, key, value, limit)
}

func (s *Service) GetDocumentByID(ctx context.Context, id int32, token string) (*Document, error) {
	return s.repo.GetDocumentByID(ctx, id, token)
}

func (s *Service) DeleteDocument(ctx context.Context, id int32, token string) error {
	return s.repo.DeleteDocument(ctx, id, token)
}
