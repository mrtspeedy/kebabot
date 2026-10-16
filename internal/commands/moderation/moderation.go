package moderation

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
	"github.com/disgoorg/snowflake/v2"
)

// Creates moderation command Definitions which can be used in commands.go
var Definitions = []discord.ApplicationCommandCreate{
	banDefiniton,
	timeoutDefinition,
	kickDefinition,
}

// Registers handlers for commands in Discord, which can also be then sent to commands.go
func Register(r handler.Router) {
	r.Command("/ban", handleBan)
	r.Command("/timeout", handleTimeout)
	r.Command("/kick", handleKick)
}

// Reply to user running a command by creating an ephemeral message that only the user can see on Discord
func reply(e *handler.CommandEvent, msg string) error {
	return e.CreateMessage(discord.NewMessageCreate().
		WithContent(msg).
		WithEphemeral(true).
		WithAllowedMentions(&discord.AllowedMentions{}))
}

// This is a quality of life function to make sure that a reason will always be given
// Returns "No Reason Provided" for an empty reason
func reasonFrom(data discord.SlashCommandInteractionData) string {
	if reason, ok := data.OptString("reason"); ok && reason != "" {
		return reason
	}
	return "No Reason Provided"
}

// Validate that the target for a moderation action is valid, by checking whether the user is in a guild
// Also checks the user's permissions too
func resolveTarget(e *handler.CommandEvent, perm discord.Permissions) (snowflake.ID, discord.User, string) {
	guildID := e.GuildID()
	if guildID == nil {
		return 0, discord.User{}, "This command can only be ran in a server."
	}

	member := e.Member()
	if member == nil || !member.Permissions.Has(perm) {
		return 0, discord.User{}, "You don't have permission to do that."
	}

	user, ok := e.SlashCommandInteractionData().OptUser("user")
	if !ok {
		return 0, discord.User{}, "Please select a user."
	}
	if user.ID == e.User().ID {
		return 0, discord.User{}, "You can't do that to yourself, dummy!"
	}
	if user.ID == e.Client().ApplicationID {
		return 0, discord.User{}, "I ain't doing nothin to myself."
	}

	return *guildID, user, ""
}
