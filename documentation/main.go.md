# cmd/kebabot/main.go

This is the file where the configuration for the bot is loaded, the connection to the database is established, and the bot is started.

## `func main()`

The main function serves a small purpose. It calls the `run` function and handles any errors produced by it.

## `func run()`

The run function actually runs the bot. If any error occurs, it returns an error to the main function.

The configuration from the `.env` is loaded into the `cfg` variable.
```go
cfg, err := config.Load()
```

The log level which slog will use to send specific logs based on the level is loaded from the `.env`.
```go
slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel})))
```

Context is loaded, which will listen for a kill signal such as CTRL+C to stop the bot.
```go
ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
```

The db url stored in the `cfg` is then checked to see if it is empty.
```go
if cfg.DatabaseURL != ""
```

The database is then connected to.
```go
pool, err := database.Connect(ctx, cfg.DatabaseURL)
```

Using the configuration that was pulled from the `.env`, the bot is created.
```go
b, err := bot.New(cfg)
```

When a kill signal is recieved, the bot waits up to 10 seconds to gracefully shut down. If it is unable to do so, the bot shuts down after the 10 seconds and shortly after, goes offline on Discord.
```go
defer func() {
    closeCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    b.Close(closeCtx)
}
```
`closeCtx` is needed here as `ctx` is already closed by the time this function is deferred.

After the bot is created, the bot is then started and commands are synced to Discord. If any error is recieved, the bot start stops.
```go
if err := b.Start(ctx); err != nil {
    return err
}
```

The channel waits for a kill signal to stop the bot, otherwise it will stay on.
```go
<-ctx.Done()
```