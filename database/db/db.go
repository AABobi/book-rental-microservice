package db

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func InitGDB() {
	dbName := "auth.db"

	_, existFileErr := os.Stat(dbName)
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	if !os.IsNotExist(existFileErr) {
		fmt.Printf("Database file %s exists\n", dbName)
	} else {
		db.AutoMigrate(&User{})
	}
	DB = db
}
