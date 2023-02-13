package context

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var Dg *discordgo.Session

// Construct a new Discord client which can be used to access the variety of Discord API functions and to set callback functions for Discord events.
func Initialize(discordToken string) {
	var err error
	Dg, err = discordgo.New("Bot " + discordToken)
	if err != nil {
		log.Fatalln("ERROR: error creating Discord session,", err)
		return
	}
}

// OpenConnection Creating a connection
func OpenConnection() {
	if err := Dg.Open(); err != nil {
		log.Fatalln("ERROR: unable to open connection,", err)
		return
	}
}
