# internal/config/config.go

This is where the configuration from the `.env` is loaded.

## `type Config struct`

This struct contains every variable from the `.env` to be read.

This is the token for the Discord bot.
```go
Token string
```

Guilds to which commands should be synced to instantly, unlike global commands which take time to sync. Stored as a slice of Discord's "Snowflake" IDs.
```go
DevGuilds []snowflake.ID
```

Checks whether to sync commands to Discord on startup.
```go
SyncCommands bool
```

Checks whether or not to enable the Message Content Intent. This will not work without enabling it on Discord first.
```go
MessageContent bool
```

This stores the database url, whose format is specified in the main [documentation](documentation.md) file.
```go
DatabaseURL string
```

This stores the log level to be used when recieving logs in the terminal.
```go
LogLevel slog.Levels
```

## `func Load()`

Start by loading in the environment variables from the `.env`.
```go
godotenv.Load()
```

Once they've been loaded, read them into the `Config` struct and name the variable `cfg`.
For the `envBool` variables, the bool after the variable indicates the default value if unprovided.
```go
cfg := &Config{
    Token:          os.Getenv("DISCORD_TOKEN"),
    SyncCommands:   envBool("SYNC_COMMANDS", true),
    MessageContent: envBool("ENABLE_MESSAGE_CONTENT", false),
    DatabaseURL:    os.Getenv("DATABASE_URL"),
}
```

After the variables have been loaded into the configuration, check if the Discord token is null, as the bot cannot connect without it.
```go
if cfg.Token == "" { return nil, errors.New("DISCORD_TOKEN is not set") }
```

If `DevGuilds` is populated, split them by commas, then parse them as Discord snowflake IDs to validate them. If any are invalid, stop the bot from starting, otherwise append them to the `DevGuilds` slice in the `cfg` variable which was defined earlier.
```go
for _, part := range strings.Split(raw, ",") {
    id, err := snowflake.Parse(strings.TrimSpace(part))
    if err != nil {
        return nil fmt.Errorf("invalid guild id %q in DEV_GUILD_IDS: %w", part, err)
    }

    cfg.DevGuilds = append(cfg.DevGuild id)
}
```

Use a switch statement to check what log level has been selected in the `.env`, falling back to Info level if an invalid level is given.
```go
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
```

Return the now loaded configuration from the `.env`.
```go
return cfg, nil
```

## `func envBool()`

The first thing that is checked is whether the `.env` key that has been passed through is unset or empty. If it is unset, then use the default value.
```go
v, ok := os.LookupEnv(key)
if !ok { return def }
```

After that, the key is converted to a boolean, which accepts values: `true`, `false`, `1`, `0`, `t`, `f`. If the key's value does not match any of those, return the default value.
```go
b, err := strconv.ParseBool(v)
if err != nil { return def }
```

If all is well then return the boolean set from the `.env`.
```go
return b
```