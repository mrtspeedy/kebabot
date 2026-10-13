# internal/bot/bot.go

This is where the functions for the bot to be created, started and closed are stored.

## `type Bot struct`

This is the client which speaks to Discord. It includes stuff like connecting to the gateway and sending messages. Basically any Discord-side part of the bot.
```go
Client *disgobot.Client
```

This holds the configuration from the `.env`. The most important part is the Discord token which will tell the bot which token to connect to. Keep in mind that bot and user tokens are different.
```go
cfg *config.Config
```

## `func New()`

An interaction router is created, which maps commands to their respective handlers.
```go
router := handler.New()
```

The command handlers are then registered to the router.
```go
commands.Register(router)
```

The bot then connects to each gateway intent specified. If you need more intents, specify them here.
```go
intents := []gateway.Intents{
    gateway.IntentGuilds,
    // More intents go here...
}
```

Event listeners are then routed through the previously initialised router, which allows the bot to listen to all the events it needs to through Discord.
```go
opts := []disgobot.ConfigOpt{
    disgobot.WithEventListeners(router),
}
```

> [!NOTE]
> The below only applies if you enable the Message Content intent on both Discord **AND** in the `.env`

If applicable, the Message Content Intent is appended to both the intents and the listener so that the bot knows it is to listen for new message creations.
```go
intents = append(intents, gateway.IntentGuildMessages, gateway.IntentMessageContent)
opts = append(opts, disgobot.WithEventListenerFunc(listeners.OnMessageCreate))
```

Then append listeners for all intents that were previously specified.
```go
opts = append(opts, disgobot.WithGatewayConfigOpts(gateway.WithIntents(intents...)))
```

The client is then built using the token and all listeners are attached to the client.
```go
client, err := disgo.New(cfg.Token, opts...)
```

## `func Start()`

If syncing commands is enabled in the `.env`, this function will first sync the commands. If any DevGuilds are specified, commands will be synced there rather than globally.
```go
handler.SyncCommands(b.Client, commands.Definitions, b.cfg.DevGuilds)
```

The gateway connection is then opened and the bot goes online on Discord.
```go
b.Client.OpenGateway(ctx)
```

## `func Close()`

This function attempts to close the gateway connection to Discord for the bot gracefully.
```go
b.Client.Close(ctx)
```