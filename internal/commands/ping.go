package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

// Define the ping command as a slash command
var pingDefinition = discord.SlashCommandCreate{
	Name:        "ping",
	Description: "Replies with pong",
}

// Handle ping by responding with "Pong!"
func handlePing(e *handler.CommandEvent) error {
	return e.CreateMessage(discord.NewMessageCreate().WithContent("Pong!"))
}
