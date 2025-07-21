package handler

import (
	"encoding/json"
	"net/http"

	"github.com/AlexanderZah/docs-management/internal/dto"
	"github.com/AlexanderZah/docs-management/internal/user"
)

type UserHandler struct {
	service *user.Service
}

func NewUserHandler(s *user.Service) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, 400, "invalid request")
		return
	}

	u, err := h.service.Register(ctx, req.AdminToken, req.Login, req.Password)
	if err != nil {
		respondWithError(w, 400, err.Error())
		return
	}

	respond(w, 200, nil, dto.RegisterResponse{Login: u.Login}, nil)
}
