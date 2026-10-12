package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var pingDefinition = discord.SlashCommandCreate{
	Name:        "ping",
	Description: "Replies with pong",
}

func handlePing(e *handler.CommandEvent) error {
	return e.CreateMessage(discord.NewMessageCreate().WithContent("Pong!"))
}
