package migration

import (
	"github.com/sitetester/ethereum-blockchain-parser/src/entity/eth"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	db.AutoMigrate(&eth.Block{})
	db.AutoMigrate(&eth.Transaction{})
	db.AutoMigrate(&eth.EventLog{})
}
