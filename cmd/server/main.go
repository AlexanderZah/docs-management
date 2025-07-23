package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/AlexanderZah/docs-management/internal/cache"
	"github.com/AlexanderZah/docs-management/internal/config"
	"github.com/AlexanderZah/docs-management/internal/docs"
	"github.com/AlexanderZah/docs-management/internal/repository"
	"github.com/AlexanderZah/docs-management/internal/router"
	"github.com/AlexanderZah/docs-management/internal/user"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	cfg := config.MustLoad()

	log.Println("cfg loaded")
	ctx := context.Background()
	pool, err := pgxpool.Connect(ctx, cfg.DbUrl)
	if err != nil {
		log.Fatalf("can't create pg pool: %s", err.Error())
	}
	userRepo := repository.NewUserRepo(pool)
	userService := user.NewService(userRepo, cfg.AdminToken)
	cacheRedis := cache.NewRedisCache(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB, 10*time.Minute)
	docsRepo := repository.NewDocsRepo(pool, cacheRedis)
	docsService := docs.NewService(docsRepo)
	r := router.NewRouter(router.Dependencies{
		UserService: userService,
		DocsService: docsService,
	})

	log.Println("Server on ", cfg.Address)
	if err := http.ListenAndServe(cfg.Address, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
