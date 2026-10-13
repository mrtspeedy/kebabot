package listeners

import (
	"log/slog"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
)

// This is an inside joke. Don't worry if you don't understand it.
const HarmonicClash = "I LOVE HARMONIC CLASH. " +
	"HARMONIC CLASH IS PEAK. " +
	"HAVE I EVER TOLD YOU THAT I LOVE HARMONIC CLASH. " +
	"HARMONIC CLASH IS SO PEAK. " +
	"LET ME TELL YOU HOW MUCH I ABSOLUTELY LOVE HARMONIC CLASH."

func OnMessageCreate(e *events.MessageCreate) {
	// Ignore all bot messages
	if e.Message.Author.Bot {
		return
	}
	// Case insensitive match on "harmonic clash"
	if !strings.Contains(strings.ToLower(e.Message.Content), "harmonic clash") {
		return
	}
	// Post a new message to the channel containing the HarmonicClash string. REST call.
	_, err := e.Client().Rest.CreateMessage(
		e.ChannelID,
		discord.NewMessageCreate().WithContent(HarmonicClash),
	)
	if err != nil {
		slog.Error("failed to glaze harmonic clash", slog.Any("err", err))
	}
}
