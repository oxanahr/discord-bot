package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type configuration struct {
	BotPrefix              string
	BotStatus              string
	BotGuildJoinMessage    string
	DiscordToken           string
	DBUser                 string
	DBPassword             string
	DBSchema               string
	DBHost                 string
	DBPort                 string
	ServerGeneralChannelID string
}

// config contains all environment variables that should be included in .env
var config *configuration

// Load loads the environment variables from the .env file
func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Error loading .env file")
	}
	//root:1234@tcp(localhost:3306)/discordBot?parseTime=true
	config = &configuration{
		BotPrefix:              os.Getenv("BOT_PREFIX"),
		BotStatus:              os.Getenv("BOT_STATUS"),
		BotGuildJoinMessage:    os.Getenv("BOT_GUILD_JOIN_MESSAGE"),
		DiscordToken:           os.Getenv("DISCORD_TOKEN"),
		DBUser:                 os.Getenv("DB_USER"),
		DBPassword:             os.Getenv("DB_PASSWORD"),
		DBSchema:               os.Getenv("DB_SCHEMA"),
		DBHost:                 os.Getenv("DB_HOST"),
		DBPort:                 os.Getenv("DB_PORT"),
		ServerGeneralChannelID: os.Getenv("SERVER_GENERAL_CHANNEL_ID"),
	}
}

func GetBotPrefix() string {
	return config.BotPrefix
}

func GetBotStatus() string {
	return config.BotStatus
}

func GetServerGeneralChannelID() string {
	return config.ServerGeneralChannelID
}

func GetDiscordToken() string {
	return config.DiscordToken
}

func GetDBUser() string {
	return config.DBUser
}

func GetDBPassword() string {
	return config.DBPassword
}

func GetDBSchema() string {
	return config.DBSchema
}

func GetDBHost() string {
	return config.DBHost
}

func GetDBPort() string {
	return config.DBPort
}
