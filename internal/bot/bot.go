package bot

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/disgoorg/disgo"
	disgobot "github.com/disgoorg/disgo/bot" // alias to disgobot so that it doesn't conflict with our bot
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/handler"

	"github.com/mrtspeedy/kebabot/internal/commands"
	"github.com/mrtspeedy/kebabot/internal/config"
	"github.com/mrtspeedy/kebabot/internal/listeners"
)

type Bot struct {
	Client *disgobot.Client // Client which speaks to Discord (Gateway Connection, Sending Messages, etc)
	cfg    *config.Config   // Holds configuration from .env
}

func New(cfg *config.Config) (*Bot, error) {
	router := handler.New()   // Interaction router, maps interactions such as /ping to their respective functions
	commands.Register(router) // Registers command handlers

	// Discord Gateway Intents (Unrequired for Slash Commands)
	intents := []gateway.Intents{
		gateway.IntentGuilds, // Guild Events such as which guilds the bot is in and any changes made in those guilds
	}
	opts := []disgobot.ConfigOpt{
		// Route interactions through earlier initialised router
		disgobot.WithEventListeners(router),
		disgobot.WithEventListenerFunc(func(*events.Ready) {
			slog.Info("bot is ready")
		}),
	}

	// If Message Content Intent enabled in .env
	if cfg.MessageContent {
		// Append to intents slice
		intents = append(intents, gateway.IntentGuildMessages, gateway.IntentMessageContent)
		// Add message creation listener to functional options, which will listen for messages
		opts = append(opts, disgobot.WithEventListenerFunc(listeners.OnMessageCreate))
	}

	// Append all intents from intents slice to functional options, which will tell the bot which events to listen for
	opts = append(opts, disgobot.WithGatewayConfigOpts(gateway.WithIntents(intents...)))

	// Builds the client using the token and the functional options
	client, err := disgo.New(cfg.Token, opts...)
	if err != nil {
		return nil, fmt.Errorf("creating disgo client: %w", err)
	}

	// Return the created bot with its client and config
	return &Bot{Client: client, cfg: cfg}, nil
}

func (b *Bot) Start(ctx context.Context) error {
	if b.cfg.SyncCommands {
		slog.Info("syncing commands", slog.Any("guilds", b.cfg.DevGuilds))
		// Syncs commands to Discord, either globally or to the dev guilds if they are populated
		if err := handler.SyncCommands(b.Client, commands.Definitions, b.cfg.DevGuilds); err != nil {
			return fmt.Errorf("syncing commands: %w", err)
		}
	}

	// Opens gateway connection and bot goes online
	if err := b.Client.OpenGateway(ctx); err != nil {
		return fmt.Errorf("opening gateway: %w", err)
	}
	return nil
}

func (b *Bot) Close(ctx context.Context) {
	// Closes the gateway connection gracefully
	b.Client.Close(ctx)
}
