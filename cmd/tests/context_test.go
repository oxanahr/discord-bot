package tests

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"testing"
)

var (
	dg    *discordgo.Session // Stores a global discordgo user session
	dgBot *discordgo.Session // Stores a global discordgo bot session
	// Token to use when authenticating the bot account
	envBotToken = os.Getenv("DISCORD_TOKEN")
)

func TestInitialize(m *testing.M) {
	fmt.Println("Init is being called.")
	if envBotToken != "" {
		if d, err := discordgo.New(envBotToken); err == nil {
			dgBot = d
		}
	}

	os.Exit(m.Run())
}
