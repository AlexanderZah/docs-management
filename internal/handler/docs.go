package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/AlexanderZah/docs-management/internal/docs"
	"github.com/AlexanderZah/docs-management/internal/dto"
)

type DocsHandler struct {
	service *docs.Service
}

func NewDocsHandler(s *docs.Service) *DocsHandler {
	return &DocsHandler{service: s}
}

func (h *DocsHandler) Upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		http.Error(w, "Invalid multipart form", http.StatusBadRequest)
		return
	}

	metaStr := r.FormValue("meta")
	var meta dto.Meta
	if err := json.Unmarshal([]byte(metaStr), &meta); err != nil {
		http.Error(w, "Invalid meta JSON", http.StatusBadRequest)
		return
	}

	jsonStr := r.FormValue("json")
	var jsonData map[string]interface{}
	if jsonStr != "" {
		if err := json.Unmarshal([]byte(jsonStr), &jsonData); err != nil {
			http.Error(w, "Invalid json field", http.StatusBadRequest)
			return
		}
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File not provided", http.StatusBadRequest)
		return
	}
	defer file.Close()
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Can't read file", http.StatusInternalServerError)
		return
	}

	createdAt := time.Now()
	doc := &docs.Document{
		Name:      meta.Name,
		IsFile:    meta.File,
		Public:    meta.Public,
		Token:     meta.Token,
		Mime:      meta.Mime,
		Grants:    meta.Grant,
		Json:      jsonData,
		Content:   fileBytes,
		CreatedAt: createdAt,
	}

	if err := h.service.UploadDocument(r.Context(), doc); err != nil {
		http.Error(w, "Failed to save document", http.StatusInternalServerError)
		return
	}

	respond(w, 200, dto.UploadDocResponse{
		Json: doc.Json,
		File: meta.Name,
	}, nil, nil)
	// // Ответ

	// json.NewEncoder(w).Encode(map[string]interface{}{
	// 	"data": map[string]interface{}{
	// 		"json": doc.Json,
	// 		"file": header.Filename,
	// 	},
	// })
}
