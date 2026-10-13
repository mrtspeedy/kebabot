package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mrtspeedy/kebabot/internal/bot"
	"github.com/mrtspeedy/kebabot/internal/config"
	"github.com/mrtspeedy/kebabot/internal/database"
)

// The main function only calls the run function, and logs a fatal error if the bot cannot be run
func main() {
	if err := run(); err != nil {
		slog.Error("fatal", slog.Any("err", err))
		os.Exit(1)
	}
}

func run() error {
	// Read the env values into a struct
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	// Logs will only log based on the log level selected in env
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel})))

	// ctx stops after recieving a kill signal such as CTRL+C
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// If the database url in the env is not empty, try to connect to it if possible
	if cfg.DatabaseURL != "" {
		pool, err := database.Connect(ctx, cfg.DatabaseURL)
		if err != nil {
			return err
		}
		// On shutdown, close connection to database
		defer pool.Close()
		slog.Info("connected to database")
	}

	// Create the bot using the discord token stored in the config
	b, err := bot.New(cfg)
	if err != nil {
		return err
	}
	defer func() {
		// Wait for bot to gracefully shutdown, else after 10 seconds, exit anyways
		closeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		b.Close(closeCtx)
	}()

	// Start the bot and sync commands
	if err := b.Start(ctx); err != nil {
		return err
	}

	slog.Info("bot running, press CTRL+C to exit")
	// Block here until shutdown signal arrives
	<-ctx.Done()
	slog.Info("shutting down")
	return nil
}
