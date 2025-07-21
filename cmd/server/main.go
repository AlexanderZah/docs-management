package main

import (
	"log/slog"
	"os"

	"github.com/AlexanderZah/docs-management/internal/config"
)

func main() {
	cfg := config.MustLoad()
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	log.Info("cfg loaded")
	_ = cfg
}
