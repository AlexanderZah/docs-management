package handler

import (
	"encoding/json"
	"net/http"

	"github.com/AlexanderZah/docs-management/internal/dto"
	"github.com/AlexanderZah/docs-management/internal/user"
)

type AuthHandler struct {
	service *user.Service
}

func NewAuthHandler(s *user.Service) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req dto.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, 400, "invalid request")
		return
	}

	token, err := h.service.Login(ctx, req.Login, req.Password)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	respond(w, 200, nil, dto.AuthResponse{Token: token}, nil)
}
