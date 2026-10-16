package moderation

import (
	"fmt"
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/omit"
)

// Defines the kick command
var kickDefinition = discord.SlashCommandCreate{
	Name:        "kick",
	Description: "Kick a member from the server",
	// The user must have the Kick Members permission to see this command
	DefaultMemberPermissions: omit.NewPtr(discord.PermissionKickMembers),
	// A user is required to ban, and a reason for the ban is optional
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionUser{
			Name:        "user",
			Description: "Member to kick",
			Required:    true,
		},
		discord.ApplicationCommandOptionString{
			Name:        "reason",
			Description: "To show in audit log",
		},
	},
}

// Handles kick command execution
func handleKick(e *handler.CommandEvent) error {
	// Check whether the target is valid
	guildID, user, msg := resolveTarget(e, discord.PermissionKickMembers)
	if msg != "" {
		return reply(e, msg)
	}

	// Get the reason from the interaction data, and return "No Reason Provided", if no reason is provided
	reason := reasonFrom(e.SlashCommandInteractionData())

	// Try to kick the user from the server
	if err := e.Client().Rest.RemoveMember(guildID, user.ID, rest.WithReason(reason)); err != nil {
		slog.Error("kick failed", slog.Any("err", err), slog.String("user", user.ID.String()))
		return reply(e, "Unable to kick selected member. Make sure to check my permissions, and that my role is above theirs.")
	}

	return reply(e, fmt.Sprintf("Kicked %s. Reason %s", user.Username, reason))
}
