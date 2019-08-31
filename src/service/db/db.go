package db

import (
	"github.com/jinzhu/gorm"
)

func GetDb() *gorm.DB {
	db, err := gorm.Open("sqlite3", "./db/blockchain.db")
	if err != nil {
		panic("Failed to connect database")
	}

	return db
}
