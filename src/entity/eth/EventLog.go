package eth

import "github.com/jinzhu/gorm"

type EventLog struct {
	gorm.Model

	Address          string // set to unique and not null
	BlockHash        string
	BlockNumber      string
	Data             string `gorm:"type:text;"`
	LogIndex         string
	Removed          bool
	TopicsString     string   `gorm:"type:text;"`
	Topics           []string `gorm:"-"`
	TransactionHash  string
	TransactionIndex string
	Date             string
}
