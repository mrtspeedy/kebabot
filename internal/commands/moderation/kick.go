package moderation

import (
	"fmt"
	"log/slog"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/omit"
)

var kickDefinition = discord.SlashCommandCreate{
	Name:                     "kick",
	Description:              "Kick a member from the server",
	DefaultMemberPermissions: omit.NewPtr(discord.PermissionKickMembers),
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

func handleKick(e *handler.CommandEvent) error {
	guildID, user, msg := resolveTarget(e, discord.PermissionKickMembers)
	if msg != "" {
		return reply(e, msg)
	}

	reason := reasonFrom(e.SlashCommandInteractionData())

	if err := e.Client().Rest.RemoveMember(guildID, user.ID, rest.WithReason(reason)); err != nil {
		slog.Error("kick failed", slog.Any("err", err), slog.String("user", user.ID.String()))
		return reply(e, "Unable to kick selected member. Make sure to check my permissions, and that my role is above theirs.")
	}

	return reply(e, fmt.Sprintf("Kicked %s. Reason %s", user.Username, reason))
}
