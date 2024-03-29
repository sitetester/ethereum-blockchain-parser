package eth

import (
	"github.com/jinzhu/gorm"
)

type Block struct {
	gorm.Model
	Number            string `gorm:"unique;not null"` // set to unique and not null
	Hash              string
	Difficulty        string
	ExtraData         string `gorm:"type:text;"`
	GasLimit          string
	GasUsed           string

	LogsBloom         string
	Miner             string
	MixHash           string
	Nonce             string
	ParentHash        string
	ReceiptsRoot      string
	Sha3Uncles        string
	Size              string
	StateRoot         string
	Timestamp         string
	TotalDifficulty   string
	Transactions      []Transaction
	EventLogs         []EventLog
	TransactionsCount int
}

