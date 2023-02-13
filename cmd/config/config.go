package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type configuration struct {
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
		log.Fatalln(" ERROR: Error loading .env file")
	}
	config = &configuration{
		DiscordToken:           os.Getenv("DISCORD_TOKEN"),
		DBUser:                 os.Getenv("DB_USER"),
		DBPassword:             os.Getenv("DB_PASSWORD"),
		DBSchema:               os.Getenv("DB_SCHEMA"),
		DBHost:                 os.Getenv("DB_HOST"),
		DBPort:                 os.Getenv("DB_PORT"),
		ServerGeneralChannelID: os.Getenv("SERVER_GENERAL_CHANNEL_ID"),
	}
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
