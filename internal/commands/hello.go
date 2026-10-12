package commands

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var helloDefinition = discord.SlashCommandCreate{
	Name:        "hello",
	Description: "Greets specified person",
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:        "name",
			Description: "Who to greet",
			Required:    true,
		},
	},
}

func handleHello(e *handler.CommandEvent) error {
	data := e.SlashCommandInteractionData()
	name, _ := data.OptString("name")

	return e.CreateMessage(discord.NewMessageCreate().
		WithContent(fmt.Sprintf("Hello, %s!", name)).
		WithAllowedMentions(&discord.AllowedMentions{}))
}
