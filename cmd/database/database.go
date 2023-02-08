package database

import (
	"fmt"
	"github.com/oxanahr/discord-bot/cmd/config"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error
	for i := 0; i < 4; i++ {
		DB, err = gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", config.GetDBUser(), config.GetDBPassword(), config.GetDBHost(), config.GetDBPort(), config.GetDBSchema())), &gorm.Config{})
		if err != nil {
			log.Println("INFO: Could not connect with the database, retrying in 5 seconds")
			time.Sleep(time.Second * 5)
		}
	}

	if err != nil {
		log.Fatalln("ERR: Could not connect with the database!")
	}
}
