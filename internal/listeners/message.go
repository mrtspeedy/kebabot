package listeners

import (
	"log/slog"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

func OnMessageCreate(e *events.MessageCreate) {
	if e.Message.Author.Bot {
		return
	}
	if !strings.EqualFold(strings.TrimSpace(e.Message.Content), "ping") {
		return
	}
	_, err := e.Client().Rest.CreateMessage(
		e.ChannelID,
		discord.NewMessageCreate().WithContent("Pong!"),
	)
	if err != nil {
		slog.Error("failed to send pong", slog.Any("err", err))
	}
}
