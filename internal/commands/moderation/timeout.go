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

const maxTimeout = 28 * duration.Day

var timeoutDefinition = discord.SlashCommandCreate{
	Name:                     "timeout",
	Description:              "Temporarily stops a member from chatting, reacting, and joining voice channels",
	DefaultMemberPermissions: omit.NewPtr(discord.PermissionModerateMembers),
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

func handleTimeout(e *handler.CommandEvent) error {
	guildID, user, msg := resolveTarget(e, discord.PermissionModerateMembers)
	if msg != "" {
		return reply(e, msg)
	}

	data := e.SlashCommandInteractionData()
	raw, _ := data.OptString("time")
	length, err := duration.ParseDuration(raw)
	if length < 1 || length > maxTimeout {
		return reply(e, fmt.Sprintf("Duration must be between 1 and up to %s", duration.Format(maxTimeout)))
	}
	reason := reasonFrom(data)

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
