package commands

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

// Defines a calculator
var calculatorDefinition = discord.SlashCommandCreate{
	Name:        "calculate",
	Description: "Basic calculator",
	// Options of number 1, operator, and number 2
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionFloat{
			Name:        "num1",
			Description: "First number",
			Required:    true,
		},
		// An operator is a restricted choice between add, subtract, multiply and divide
		discord.ApplicationCommandOptionString{
			Name:        "operator",
			Description: "Operator to perform calculation",
			Required:    true,
			Choices: []discord.ApplicationCommandOptionChoiceString{
				{Name: "Add", Value: "+"},
				{Name: "Subtract", Value: "-"},
				{Name: "Multiply", Value: "*"},
				{Name: "Divide", Value: "/"},
			},
		},
		discord.ApplicationCommandOptionFloat{
			Name:        "num2",
			Description: "Second number",
			Required:    true,
		},
	},
}

// Handles calculator execution
func handleCalculator(e *handler.CommandEvent) error {
	data := e.SlashCommandInteractionData()
	num1, _ := data.OptFloat("num1")
	operator, _ := data.OptString("operator")
	num2, _ := data.OptFloat("num2")
	var result float64

	// Switch case for each possible operator
	switch operator {
	case "+":
		result = num1 + num2
	case "-":
		result = num1 - num2
	case "*":
		result = num1 * num2
	case "/":
		// Handle divison by zero
		if num2 == 0 {
			return e.CreateMessage(discord.NewMessageCreate().WithContent("Division by zero is not allowed!"))
		}
	default:
		// Return an error if it appears the operator is invalid
		return e.CreateMessage(discord.NewMessageCreate().WithContent("Something went wrong. Check your inputs."))
	}

	return e.CreateMessage(discord.NewMessageCreate().
		WithContent(fmt.Sprintf("%v %v %v is %v!", num1, operator, num2, result)).
		// Make sure the bot only mentions those whom the user has permissions to mention
		WithAllowedMentions(&discord.AllowedMentions{}))
}
