package handler

import (
	"encoding/json"
	"net/http"

	"github.com/AlexanderZah/docs-management/internal/auth"
	"github.com/AlexanderZah/docs-management/internal/dto"
)

type AuthHandler struct {
	service *auth.Service
}

func NewAuthHandler(s *auth.Service) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, 400, "invalid request")
		return
	}

	token, err := h.service.Authenticate(req.Login, req.Password)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}

	respond(w, 200, nil, dto.AuthResponse{Token: token}, nil)
}
