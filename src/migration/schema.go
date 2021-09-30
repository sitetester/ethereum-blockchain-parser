package migration

import (
	"github.com/jinzhu/gorm"
	"github.com/sitetester/ethereum-blockchain-parser/src/entity/eth"
)

func RunMigrations(db *gorm.DB) {
	db.AutoMigrate(&eth.Block{})
	db.AutoMigrate(&eth.Transaction{})
	db.AutoMigrate(&eth.EventLog{})
}
