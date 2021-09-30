package eth

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/sitetester/ethereum-blockchain-parser/src/entity/eth"
	"strconv"
)

type ImportManager struct {
	bigQueryClient BigQueryClient
}

func (manager ImportManager) ManageImport(db *gorm.DB) {
	var parser BlocksParser

	var ethBlock eth.Block
	db.Limit(1).Select("Number").Order("Number desc").Find(&ethBlock)
	fmt.Printf("lastScannedBlockNumber: %s \n", ethBlock.Number)

	var blockNumber int
	blockNumberTemp, _ := strconv.Atoi(ethBlock.Number)
	blockNumber = blockNumberTemp

	parsedBlocks := parser.ParseBlocks(blockNumber + 1)
	fmt.Printf("\n All blocks parsed!!!\n")

	for _, parsedBlock := range parsedBlocks {
		var bqTransactionList []eth.BigQueryTransaction

		for i, _ := range parsedBlock.Transactions {
			transaction := &parsedBlock.Transactions[i]

			gasInt, _ := strconv.Atoi(transaction.Gas)
			gasPriceFloat, _ := strconv.ParseFloat(transaction.GasPrice, 64)

			bqTransaction := eth.BigQueryTransaction{
				Hash:             transaction.Hash,
				Sender:           transaction.From,
				Receiver:         transaction.To,
				BlockHash:        transaction.BlockHash,
				BlockNumber:      transaction.BlockNumber,
				Gas:              gasInt,
				GasPrice:         gasPriceFloat,
				Value:            transaction.Value,
				TransactionIndex: transaction.TransactionIndex,
				Nonce:            transaction.Nonce,
				Date:             transaction.Date,
			}

			bqTransactionList = append(bqTransactionList, bqTransaction)
		}

		inserted := manager.bigQueryClient.InsertRows(TransactionsTable, bqTransactionList)
		if inserted {
			db.Create(&parsedBlock)
		}
	}
}
