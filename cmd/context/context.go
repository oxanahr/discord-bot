package context

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var Dg *discordgo.Session

// Initializing discord session
func Initialize(discordToken string) {
	var err error
	Dg, err = discordgo.New("Bot " + discordToken)
	if err != nil {
		log.Fatalln("ERROR: error creating Discord session,", err)
		return
	}
}

// Creating a connection
func OpenConnection() {
	if err := Dg.Open(); err != nil {
		log.Fatalln("ERROR: unable to open connection,", err)
		return
	}
}
