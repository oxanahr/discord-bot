package main

import (
	"fmt"
	"github.com/madflojo/tasks"
	"github.com/oxanahr/discord-bot/cmd/bot"
	"github.com/oxanahr/discord-bot/cmd/config"
	"github.com/oxanahr/discord-bot/cmd/database"
	"github.com/oxanahr/discord-bot/cmd/database/migrations"
	"github.com/oxanahr/discord-bot/cmd/handlers"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Load environment variables
	config.Load()

	// Connect to database and run migrations
	database.Connect()
	migrations.AutoMigrate()

	// Start the bot
	bot.Start()

	scheduler := tasks.New()
	defer scheduler.Stop()

	_, err := scheduler.Add(&tasks.Task{
		//Interval:   24 * time.Hour,
		Interval: 30 * time.Second,
		//StartAfter: time.Now().Add(10 * time.Second), // run at X:00 every day?
		TaskFunc: func() error {
			handlers.PingUsers() //error handling
			handlers.Summary()   //error handling
			return nil
		},
	})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Bot is running. Press Ctrl + C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	bot.Stop()
}
