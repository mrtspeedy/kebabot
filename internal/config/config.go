package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/disgoorg/snowflake/v2"
	"github.com/joho/godotenv"
)

type Config struct {
	Token          string
	DevGuilds      []snowflake.ID
	SyncCommands   bool
	MessageContent bool
	DatabaseURL    string
	LogLevel       slog.Level
}

func Load() (*Config, error) {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("loading .env: %w", err)
	}

	cfg := &Config{
		Token:          os.Getenv("DISCORD_TOKEN"),
		SyncCommands:   envBool("SYNC_COMMAND", true),
		MessageContent: envBool("ENABLE_MESSAGE_CONTENT", false),
		DatabaseURL:    os.Getenv("DATABASE_URL"),
	}
	if cfg.Token == "" {
		return nil, errors.New("DISCORD_TOKEN is not set")
	}
	if raw := strings.TrimSpace(os.Getenv("DEV_GUILD_IDS")); raw != "" {
		for _, part := range strings.Split(raw, ",") {
			id, err := snowflake.Parse(strings.TrimSpace(part))
			if err != nil {
				return nil, fmt.Errorf("invalid guild id %q in DEV_GUILD_IDS: %w", part, err)
			}
			cfg.DevGuilds = append(cfg.DevGuilds, id)
		}
	}
	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "debug":
		cfg.LogLevel = slog.LevelDebug
	case "warn":
		cfg.LogLevel = slog.LevelWarn
	case "error":
		cfg.LogLevel = slog.LevelError
	default:
		cfg.LogLevel = slog.LevelInfo
	}

	return cfg, nil
}

func envBool(key string, def bool) bool {
	v, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return def
	}
	return b
}
