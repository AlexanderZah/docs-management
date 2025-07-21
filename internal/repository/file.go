package repository

import (
	"context"
	"os"
	"path/filepath"

	"github.com/AlexanderZah/docs-management/internal/docs"
)

type FileRepository struct {
	basePath string
}

func NewFileRepository(path string) *FileRepository {
	return &FileRepository{basePath: path}
}

func (r *FileRepository) SaveDocument(ctx context.Context, doc *docs.Document) error {
	// Просто сохраняем файл в базовую директорию
	filePath := filepath.Join(r.basePath, doc.Name)
	return os.WriteFile(filePath, doc.File, 0644)
}
