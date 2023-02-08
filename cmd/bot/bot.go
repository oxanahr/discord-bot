package bot

import (
	"github.com/oxanahr/discord-bot/cmd/config"
	"github.com/oxanahr/discord-bot/cmd/context"
	"github.com/oxanahr/discord-bot/cmd/handlers"
	"math/rand"
	"time"
)

func Start() {
	rand.Seed(time.Now().UnixNano())
	context.Initialize(config.GetDiscordToken())
	handlers.AddHandlers()
	context.OpenConnection()
}

func Stop() {
	context.Dg.Close()
}
