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
	r.Use(handler.NewAuthMiddleware(dep.UserService).Handle)
	userHandler := handler.NewUserHandler(dep.UserService)
	docsHandler := handler.NewDocsHandler(dep.DocsService)
	r.Route("/api", func(r chi.Router) {
		r.Post("/register", userHandler.Register)
		r.Post("/auth", userHandler.Login)
		r.Delete("/auth/{token}", userHandler.Logout)
		r.Post("/docs", docsHandler.Upload)
		r.Get("/docs", docsHandler.Get)
		r.Head("/docs", docsHandler.Get)
		r.Get("/docs/{id}", docsHandler.GetByID)
		r.Head("/docs/{id}", docsHandler.GetByID)
		r.Delete("/docs/{id}", docsHandler.Delete)
	})

	return r
}
