package eth

import (
	"fmt"
	"github.com/sitetester/ethereum-blockchain-parser/src/entity/eth"
	"gorm.io/gorm"
	"time"
)

type ImportManager struct {
	bigQueryClient BigQueryClient
}

func (manager ImportManager) ManageImport(db *gorm.DB) {
	var parser BlocksParser

	var ethBlock eth.Block
	db.Limit(1).Select("NumberInt").Order("number_int desc").Find(&ethBlock)
	fmt.Printf("lastScannedBlockNumber: %v \n", ethBlock.NumberInt)

	const size = 10
	start := ethBlock.NumberInt + 1

	for {
		parsedBlocks := parser.ParseBlocks(start, size)
		db.Create(&parsedBlocks)

		time.Sleep(5 * time.Second) // sometimes remote server stops responding
		start += size
	}
}
