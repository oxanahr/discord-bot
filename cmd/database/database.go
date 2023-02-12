package database

import (
	"fmt"
	"github.com/oxanahr/discord-bot/cmd/config"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	// Print Slow SQL and happening errors by default
	var err error
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,  // Slow SQL threshold
			LogLevel:                  logger.Error, // Log level
			IgnoreRecordNotFoundError: true,         // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,        // Disable color
		},
	)

	for i := 0; i < 4; i++ {
		DB, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.GetDBUser(), config.GetDBPassword(), config.GetDBHost(), config.GetDBPort(), config.GetDBSchema())), &gorm.Config{
			Logger: newLogger,
		})
		if err != nil {
			log.Println("INFO: Could not connect with the database, retrying in 5 seconds")
			time.Sleep(time.Second * 5)
		}
	}

	if err != nil {
		log.Fatalln("ERR: Could not connect with the database!")
	}
}
