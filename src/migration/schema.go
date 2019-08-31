package migration

import (
	"blockchain/src/entity/eth"

	"github.com/jinzhu/gorm"
)

func RunMigrations(db *gorm.DB) {
	db.AutoMigrate(&eth.Block{})
	db.AutoMigrate(&eth.Transaction{})
	db.AutoMigrate(&eth.EventLog{})
}
