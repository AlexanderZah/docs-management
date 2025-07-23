package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/AlexanderZah/docs-management/internal/docs"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const documentsTable = "documents"

type DocsRepo struct {
	db *pgxpool.Pool
}

func NewDocsRepo(db *pgxpool.Pool) *DocsRepo {
	return &DocsRepo{db: db}
}

func (r *DocsRepo) SaveDocument(ctx context.Context, doc *docs.Document) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	doc.CreatedAt = time.Now()

	var jsonData []byte
	if len(doc.Json) > 0 {
		jsonData, err = json.Marshal(doc.Json)
		if err != nil {
			return fmt.Errorf("failed to marshal json: %w", err)
		}
	}

	const query = `
        INSERT INTO documents (
            name, is_file, is_public, token, mime, grants, json_data, content
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8
        ) RETURNING id
    `

	err = tx.QueryRow(ctx, query,
		doc.Name,
		doc.IsFile,
		doc.Public,
		doc.Token,
		doc.Mime,
		doc.Grants,
		jsonData,
		doc.Content,
	).Scan(&doc.ID)

	if err != nil {
		return fmt.Errorf("failed to insert document: %w", err)
	}

	return tx.Commit(ctx)
}

func (r *DocsRepo) GetDocuments(ctx context.Context, token, login, key, value string, limit int) ([]docs.Document, error) {
	query := `
        SELECT id, name, is_file, is_public, token, mime, grants, json_data, content, created_at
        FROM documents
        WHERE (token = $1 OR $2 = ANY(grants))
    `
	args := []interface{}{token, login}


	allowedKeys := map[string]bool{
		"name":       true,
		"is_file":    true,
		"is_public":  true,
		"mime":       true,
		"json_data":  true,
		"content":    true,
		"created_at": true,
	}

	if key != "" && value != "" {
		if allowedKeys[key] {
			query += fmt.Sprintf(" AND %s = $%d", key, len(args)+1)
			args = append(args, value)
		} else {
			return nil, fmt.Errorf("invalid key for filtering")
		}
	}

	query += " ORDER BY name, created_at"
	if limit != 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query documents: %w", err)
	}
	defer rows.Close()

	var docsList []docs.Document
	for rows.Next() {
		var doc docs.Document
		var jsonData []byte
		err := rows.Scan(
			&doc.ID,
			&doc.Name,
			&doc.IsFile,
			&doc.Public,
			&doc.Token,
			&doc.Mime,
			&doc.Grants,
			&jsonData,
			&doc.Content,
			&doc.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		if len(jsonData) > 0 {
			if err := json.Unmarshal(jsonData, &doc.Json); err != nil {
				return nil, fmt.Errorf("failed to unmarshal json: %w", err)
			}
		}
		docsList = append(docsList, doc)
	}

	return docsList, nil
}

func (r *DocsRepo) GetDocumentByID(ctx context.Context, id int32, login string) (*docs.Document, error) {
	var doc docs.Document
	var jsonData []byte

	query := `
        SELECT id, name, is_file, is_public, token, mime, grants, json_data, content, created_at
        FROM documents
        WHERE id = $1 AND ($2 = ANY(grants) OR token = $2)
    `
	err := r.db.QueryRow(ctx, query, id, login).Scan(
		&doc.ID,
		&doc.Name,
		&doc.IsFile,
		&doc.Public,
		&doc.Token,
		&doc.Mime,
		&doc.Grants,
		&jsonData,
		&doc.Content,
		&doc.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	if len(jsonData) > 0 {
		if err := json.Unmarshal(jsonData, &doc.Json); err != nil {
			return nil, fmt.Errorf("failed to unmarshal json: %w", err)
		}
	}

	return &doc, nil
}

func (r *DocsRepo) DeleteDocument(ctx context.Context, id int32, token string) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	result, err := tx.Exec(ctx, `
        DELETE FROM documents
        WHERE id = $1 AND token = $2
    `, id, token)

	if err != nil {
		return fmt.Errorf("failed to delete document: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("document not found or unauthorized")
	}

	return tx.Commit(ctx)
}
