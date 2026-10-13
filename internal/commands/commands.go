package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

// Create command definitions in Discord, which tells Discord information about a command
var Definitions = []discord.ApplicationCommandCreate{
	pingDefinition,
	helloDefinition,
	calculatorDefinition,
}

// Registers the commands and their handlers, telling Discord how to handle each command
func Register(r handler.Router) {
	r.Command("/ping", handlePing)
	r.Command("/hello", handleHello)
	r.Command("/calculate", handleCalculator)
}
