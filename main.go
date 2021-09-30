package main

import (
	"github.com/jinzhu/gorm"
	"github.com/sitetester/ethereum-blockchain-parser/src/migration"
	"github.com/sitetester/ethereum-blockchain-parser/src/service/eth"
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
