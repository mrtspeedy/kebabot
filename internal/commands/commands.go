package commands

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/handler"
)

var Definitions = []discord.ApplicationCommandCreate{
	pingDefinition,
	helloDefinition,
}

func Register(r handler.Router) {
	r.Command("/ping", handlePing)
	r.Command("/hello", handleHello)
}
