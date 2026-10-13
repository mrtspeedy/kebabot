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
	Token          string         // Discord bot token
	DevGuilds      []snowflake.ID // Guilds to sync commands to; synced quickly unlike global commands (optional, mainly for testing in certain servers)
	SyncCommands   bool           // Sync commands to discord on startup
	MessageContent bool           // Whether or not the Message Content Intent should be enabled (required for message listeners)
	DatabaseURL    string         // Database to connect to, can be postgres, sqlite, or any other valid db
	LogLevel       slog.Level     // Log level to be used when recieving bot logs in logs
}

func Load() (*Config, error) {
	// Load environment variables from .env
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("loading .env: %w", err)
	}

	cfg := &Config{
		Token:          os.Getenv("DISCORD_TOKEN"),               // Discord token from the env variables
		SyncCommands:   envBool("SYNC_COMMAND", true),            // Whether or not to sync commands on startup, default true
		MessageContent: envBool("ENABLE_MESSAGE_CONTENT", false), // Whether or not to enable Message Content Intent, default false
		DatabaseURL:    os.Getenv("DATABASE_URL"),                // Database url from env
	}
	// Check whether the Discord token is null
	if cfg.Token == "" {
		return nil, errors.New("DISCORD_TOKEN is not set")
	}
	// Check whether the dev guilds are populated
	if raw := strings.TrimSpace(os.Getenv("DEV_GUILD_IDS")); raw != "" {
		// If so, split them by commas, put the result in id, and return an error if an invalid ID is detected
		for _, part := range strings.Split(raw, ",") {
			id, err := snowflake.Parse(strings.TrimSpace(part))
			if err != nil {
				return nil, fmt.Errorf("invalid guild id %q in DEV_GUILD_IDS: %w", part, err)
			}
			// Append each ID to the dev guilds slice
			cfg.DevGuilds = append(cfg.DevGuilds, id)
		}
	}
	// Get log level from env, default to Info if invalid
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
	// Distinguishes between unset and empty
	v, ok := os.LookupEnv(key)
	// If unset, then return the default value
	if !ok {
		return def
	}
	// Convert to a bool
	// Accepts true, false, 1, 0, t, f
	b, err := strconv.ParseBool(v)
	// If invalid, then return the default value
	if err != nil {
		return def
	}
	return b
}
