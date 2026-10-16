package moderation

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/omit"
)

// Defines the ban command
var banDefiniton = discord.SlashCommandCreate{
	Name:                     "ban",
	Description:              "Ban a member from the server",
	// The user must have the Ban Members permission to see this command
	DefaultMemberPermissions: omit.NewPtr(discord.PermissionBanMembers),
	// A user is required to ban, plus the reason for the ban and the messages that should be deleted are optional
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionUser{
			Name:        "user",
			Description: "Member to ban",
			Required:    true,
		},
		discord.ApplicationCommandOptionString{
			Name:        "reason",
			Description: "To show in the audit log",
		},
		discord.ApplicationCommandOptionInt{
			Name:        "delete_days",
			Description: "Delete messages from user, from the past entered number of days",
		},
	},
}

// Handles ban command execution
func handleBan(e *handler.CommandEvent) error {
	// Check whether the target is valid
	guildID, user, msg := resolveTarget(e, discord.PermissionBanMembers)
	if msg != "" {
		return reply(e, msg)
	}

	// Get the data from the interaction
	data := e.SlashCommandInteractionData()

	// Get the reason from the interaction, returning "No Reason Provided" when none is provided
	reason := reasonFrom(data)

	days, _ := data.OptInt("delete_days")
	days = min(max(days, 0), 7) // Only delete their messages from the past 0 to 7 days

	// Try to ban the selected member, and delete their messages from the past 'days' number of days
	err := e.Client().Rest.AddBan(
		guildID, user.ID,
		time.Duration(days)*24*time.Hour,
		rest.WithReason(reason),
	)
	if err != nil {
		slog.Error("ban failed", slog.Any("err", err), slog.String("user", user.ID.String()))
		return reply(e, "Unable to ban selected member. Make sure to check my permissions, and that my role is above theirs.")
	}

	return reply(e, fmt.Sprintf("Banned %s. Reason: %s", user.Username, reason))
}
