package bot

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/disgoorg/disgo"
	disgobot "github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"

	"github.com/mrtspeedy/kebabot/internal/commands"
	"github.com/mrtspeedy/kebabot/internal/config"
	"github.com/mrtspeedy/kebabot/internal/listeners"
)

type Bot struct {
	Client *disgobot.Client
	cfg    *config.Config
}

func New(cfg *config.Config) (*Bot, error) {
	router := handler.New()
	commands.Register(router)

	intents := []gateway.Intents{gateway.IntentGuilds}
	opts := []disgobot.ConfigOpt{
		disgobot.WithEventListeners(router),
		disgobot.WithEventListenerFunc(func(*events.Ready) {
			slog.Info("bot is ready")
		}),
	}

	if cfg.MessageContent {
		intents = append(intents, gateway.IntentGuildMessages, gateway.IntentMessageContent)
		opts = append(opts, disgobot.WithEventListenerFunc(listeners.OnMessageCreate))
	}

	opts = append(opts, disgobot.WithGatewayConfigOpts(gateway.WithIntents(intents...)))

	client, err := disgo.New(cfg.Token, opts...)
	if err != nil {
		return nil, fmt.Errorf("creating disgo client: %w", err)
	}

	return &Bot{Client: client, cfg: cfg}, nil
}

func (b *Bot) Start(ctx context.Context) error {
	if b.cfg.SyncCommands {
		slog.Info("syncing commands", slog.Any("guilds", b.cfg.DevGuilds))
		if err := handler.SyncCommands(b.Client, commands.Definitions, b.cfg.DevGuilds); err != nil {
			return fmt.Errorf("syncing commands: %w", err)
		}
	}

	if err := b.Client.OpenGateway(ctx); err != nil {
		return fmt.Errorf("opening gateway: %w", err)
	}
	return nil
}

func (b *Bot) Close(ctx context.Context) {
	b.Client.Close(ctx)
}
