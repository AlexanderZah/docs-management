package router

import (
	"net/http"

	"github.com/AlexanderZah/docs-management/internal/docs"
	"github.com/AlexanderZah/docs-management/internal/handler"
	"github.com/AlexanderZah/docs-management/internal/user"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type Dependencies struct {
	UserService *user.Service
	DocsService *docs.Service
}

func NewRouter(dep Dependencies) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	userHandler := handler.NewUserHandler(dep.UserService)
	authHandler := handler.NewAuthHandler(dep.UserService)
	docsHandler := handler.NewDocsHandler(dep.DocsService)
	r.Route("/api", func(r chi.Router) {
		r.Post("/register", userHandler.Register)
		r.Post("/auth", authHandler.Login)
		r.Post("/docs", docsHandler.Upload)
	})

	return r
}
