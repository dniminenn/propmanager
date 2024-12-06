package db

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("propmanager.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}
