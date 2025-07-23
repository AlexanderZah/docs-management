package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/AlexanderZah/docs-management/internal/docs"
	"github.com/AlexanderZah/docs-management/internal/dto"
	"github.com/go-chi/chi/v5"
)

type DocsHandler struct {
	Service *docs.Service
}

func NewDocsHandler(s *docs.Service) *DocsHandler {
	return &DocsHandler{Service: s}
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

	if err := h.Service.UploadDocument(r.Context(), doc); err != nil {
		http.Error(w, "Failed to save document", http.StatusInternalServerError)
		return
	}

	respond(w, 200, dto.UploadDocResponse{
		Json: doc.Json,
		File: meta.Name,
	}, nil, nil)
}

func (h *DocsHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token := r.URL.Query().Get("token")
	if token == "" {
		token = r.Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")
	}
	if token == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}

	login := r.URL.Query().Get("login")
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")
	limitStr := r.URL.Query().Get("limit")

	limit := 100
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	docs, err := h.Service.GetDocuments(ctx, token, login, key, value, limit)
	if err != nil {
		http.Error(w, "failed to fetch documents", http.StatusInternalServerError)
		return
	}
	var docsResp []dto.GetDocResponse
	for _, d := range docs {
		docsResp = append(docsResp, dto.GetDocResponse{
			ID:        d.ID,
			Name:      d.Name,
			Mime:      d.Mime,
			IsFile:    d.IsFile,
			Public:    d.Public,
			CreatedAt: d.CreatedAt.Format("2006-01-02 15:04:05"),
			Grants:    d.Grants,
		})
	}

	data := map[string]interface{}{
		"docs": docsResp,
	}
	respond(w, 200, data, nil, nil)
}

func (h *DocsHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	token := r.URL.Query().Get("token")
	if token == "" {
		token = r.Header.Get("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")
	}
	if token == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}

	doc, err := h.Service.GetDocumentByID(ctx, int32(id), token)
	if err != nil {
		http.Error(w, "document not found or access denied", http.StatusNotFound)
		return
	}

	if doc.IsFile {
		w.Header().Set("Content-Type", doc.Mime)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(doc.Content)
		return
	}
	respond(w, 200, doc.Json, nil, nil)
}

func (h *DocsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		http.Error(w, "Invalid document ID", http.StatusBadRequest)
		return
	}

	token := r.URL.Query().Get("token")
	if token == "" {
		token = strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
	}
	if token == "" {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}

	err = h.Service.DeleteDocument(ctx, int32(id), token)
	if err != nil {
		http.Error(w, "Failed to delete document or not authorized", http.StatusForbidden)
		return
	}

	resp := map[string]bool{
		idParam: true,
	}
	respond(w, 200, nil, resp, nil)

}
