package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mrtspeedy/kebabot/internal/bot"
	"github.com/mrtspeedy/kebabot/internal/config"
	"github.com/mrtspeedy/kebabot/internal/database"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel})))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if cfg.DatabaseURL != "" {
		pool, err := database.Connect(ctx, cfg.DatabaseURL)
		if err != nil {
			return err
		}
		defer pool.Close()
		slog.Info("connected to database")
	}

	b, err := bot.New(cfg)
	if err != nil {
		return err
	}
	defer func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		b.Close(closeCtx)
	}()

	if err := b.Start(ctx); err != nil {
		return err
	}

	slog.Info("bot running, press CTRL+C to exit")
	<-ctx.Done()
	slog.Info("shutting down")
	return nil
}
