package migrations

import (
	"github.com/oxanahr/discord-bot/cmd/database"
	"github.com/oxanahr/discord-bot/cmd/models"
)

// AutoMigrate automatically migrates models to database tables
func AutoMigrate() {
	database.DB.AutoMigrate(models.Task{})
	database.DB.AutoMigrate(models.Comment{})
}
