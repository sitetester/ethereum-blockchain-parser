package main

import (
	"blockchain/src/migration"
	"blockchain/src/service/eth"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	db, err := gorm.Open("sqlite3", "./db/blockchain.db")
	if err != nil {
		panic("Failed to connect database!")
	}
	defer db.Close()

	migration.RunMigrations(db)

	var importManager eth.ImportManager
	importManager.ManageImport(db)
}
