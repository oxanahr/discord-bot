package bot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/oxanahr/discord-bot/cmd/config"
	"github.com/oxanahr/discord-bot/cmd/context"
	"github.com/oxanahr/discord-bot/cmd/handlers"
	"log"
	"math/rand"
	"time"
)

var (
	registeredCommands []*discordgo.ApplicationCommand
	GuildID            = ""
)

// Opens a session to Discord API and registers all executable commands
func Start() {
	rand.Seed(time.Now().UnixNano())
	context.Initialize(config.GetDiscordToken())
	context.OpenConnection()
	handlers.ReadyHandler()
	handlers.RegisterCommands()
	handlers.MessageCreateHandler()
}

// Closes session with Discord API and deletes registered commands
func Stop() {
	for _, v := range registeredCommands {
		err := context.Dg.ApplicationCommandDelete(context.Dg.State.User.ID, GuildID, v.ID)
		if err != nil {
			log.Printf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}

	context.Dg.Close()
}
