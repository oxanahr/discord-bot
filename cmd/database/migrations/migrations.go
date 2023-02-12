package migrations

import (
	"github.com/oxanahr/discord-bot/cmd/database"
	"github.com/oxanahr/discord-bot/cmd/models"
)

func AutoMigrate() {
	database.DB.AutoMigrate(models.Task{})
	database.DB.AutoMigrate(models.Comment{})
}
