package commands

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

// Define hello as a command that takes in a name as a parameter
var helloDefinition = discord.SlashCommandCreate{
	Name:        "hello",
	Description: "Greets specified person",
	// Required name parameter shown in Discord
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionString{
			Name:        "name",
			Description: "Who to greet",
			Required:    true,
		},
	},
}

// Handles hello by replying with hello and the entered parameter name
func handleHello(e *handler.CommandEvent) error {
	data := e.SlashCommandInteractionData() // Command data from interaction, including parameters
	name, _ := data.OptString("name")       // Read the name parameter from interaction, ok is ignored as parameter is required

	return e.CreateMessage(discord.NewMessageCreate().
		WithContent(fmt.Sprintf("Hello, %s!", name)).
		// Allowed mentions is empty so no one can be mentioned in the message
		WithAllowedMentions(&discord.AllowedMentions{}))
}
