package moderation

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/omit"
	"github.com/mrtspeedy/kebabot/internal/duration"
)

// Sets the maximum timout to 28 days
const maxTimeout = 28 * duration.Day

// Defines the timeout command
var timeoutDefinition = discord.SlashCommandCreate{
	Name:        "timeout",
	Description: "Temporarily stops a member from chatting, reacting, and joining voice channels",
	// The user must have the Moderate Members permission to see this command
	DefaultMemberPermissions: omit.NewPtr(discord.PermissionModerateMembers),
	// A user is required to time out, a time is required and a reason is optional
	Options: []discord.ApplicationCommandOption{
		discord.ApplicationCommandOptionUser{
			Name:        "user",
			Description: "Member to time out",
			Required:    true,
		},
		discord.ApplicationCommandOptionString{
			Name:        "time",
			Description: "e.g. 30m, 2h, 1d12h (s, m, h, d). Max 28 days.",
			Required:    true,
		},
		discord.ApplicationCommandOptionString{
			Name:        "reason",
			Description: "To show in audit log",
		},
	},
}

// Handles timeout command execution
func handleTimeout(e *handler.CommandEvent) error {
	// Check whether the target is valid
	guildID, user, msg := resolveTarget(e, discord.PermissionModerateMembers)
	if msg != "" {
		return reply(e, msg)
	}

	// Get the data from the interaction
	data := e.SlashCommandInteractionData()

	// Get the raw time entered by the user
	raw, _ := data.OptString("time")

	// Parse the time into one the program can understand
	length, err := duration.ParseDuration(raw)

	// Check if the time given is valid
	if length < 1 || length > maxTimeout {
		return reply(e, fmt.Sprintf("Duration must be between 1m and up to %s", duration.Format(maxTimeout)))
	}

	// Get the reason from the interaction data
	reason := reasonFrom(data)

	// Try to time out the selected user until a specific time
	until := time.Now().Add(length)
	_, err = e.Client().Rest.UpdateMember(
		guildID, user.ID,
		discord.MemberUpdate{CommunicationDisabledUntil: omit.NewPtr(until)},
		rest.WithReason(reason),
	)
	if err != nil {
		slog.Error("timeout failed", slog.Any("err", err), slog.String("user", user.ID.String()))
		return reply(e, "Unable to time out selected member. Make sure to check my permissions, and that my role is above theirs.")
	}

	return reply(e, fmt.Sprintf("Timed out %s for %s. Reason: %s", user.Username, duration.Format(length), reason))
}
