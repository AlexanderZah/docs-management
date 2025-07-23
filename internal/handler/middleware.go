package handler

import (
	"net/http"
	"strings"

	"github.com/AlexanderZah/docs-management/internal/user"
)

type AuthMiddleware struct {
	service *user.Service
}

func NewAuthMiddleware(s *user.Service) *AuthMiddleware {
	return &AuthMiddleware{service: s}
}

func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if strings.HasPrefix(r.URL.Path, "/api/auth/") || r.URL.Path == "/api/auth" || r.URL.Path == "/api/register" {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "invalid Authorization format", http.StatusUnauthorized)
			return
		}
		token := parts[1]

		_, err := m.service.GetSessionByToken(r.Context(), token)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
