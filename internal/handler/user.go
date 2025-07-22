package handler

import (
	"encoding/json"
	"net/http"

	"github.com/AlexanderZah/docs-management/internal/dto"
	"github.com/AlexanderZah/docs-management/internal/user"
	"github.com/go-chi/chi/v5"
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

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
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

	respond(w, 200, nil, dto.LoginResponse{Token: token}, nil)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := chi.URLParam(r, "token")

	if token == "" {
		respondWithError(w, 400, "missing token")
		return
	}

	err := h.service.Logout(ctx, token)
	if err != nil {
		respondWithError(w, 500, "failed to logout")
		return
	}

	respond(w, 200, nil, dto.LogoutResponse{
		Response: map[string]bool{
			token: true,
		},
	}, nil)
}
