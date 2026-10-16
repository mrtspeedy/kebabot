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

var banDefiniton = discord.SlashCommandCreate{
	Name:                     "ban",
	Description:              "Ban a member from the server",
	DefaultMemberPermissions: omit.NewPtr(discord.PermissionBanMembers),
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

func handleBan(e *handler.CommandEvent) error {
	guildID, user, msg := resolveTarget(e, discord.PermissionBanMembers)
	if msg != "" {
		return reply(e, msg)
	}

	data := e.SlashCommandInteractionData()
	reason := reasonFrom(data)

	days, _ := data.OptInt("delete_days")
	days = min(max(days, 0), 7) // Only delete their messages from the past 0 to 7 days

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
