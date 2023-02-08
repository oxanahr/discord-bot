package migrations

import (
	"github/oxanahr/discordBot/cmd/database"
	"github/oxanahr/discordBot/cmd/models"
)

func AutoMigrate() {
	database.DB.AutoMigrate(models.Task{})
}
