package config

import (
	"testing"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	config := &configuration{
		DiscordToken:           "hajj78398393-iiII",
		DBPort:                 "3306",
		DBHost:                 "localhost",
		DBSchema:               "dbSchema",
		DBUser:                 "user",
		DBPassword:             "password",
		ServerGeneralChannelID: "200303992711",
	}

	if config.DiscordToken != "hajj78398393-iiII" {
		t.Errorf("config.DiscordToken == %v", config.DiscordToken)
	}
	if config.DBPort != "3306" {
		t.Errorf("config.DBPort == %v", config.DBPort)
	}
	if config.DBHost != "localhost" {
		t.Errorf("config.DBHost == %v", config.DBHost)
	}
	if config.DBSchema != "dbSchema" {
		t.Errorf("config.DBSchema == %v", config.DBSchema)
	}
	if config.DBUser != "user" {
		t.Errorf("config.DBUser == %v", config.DBUser)
	}
	if config.DBPassword != "password" {
		t.Errorf("config.DBPassword == %v", config.DBPassword)
	}
	if config.ServerGeneralChannelID != "200303992711" {
		t.Errorf("config.ServerGeneralChannelID == %v", config.ServerGeneralChannelID)
	}
}
