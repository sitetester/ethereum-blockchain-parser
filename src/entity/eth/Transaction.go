package eth

type Transaction struct {
	Hash             string `gorm:"unique;not null"` // set to unique and not null
	BlockNumber      string
	BlockHash        string
	From             string
	To               string
	Gas              string
	GasPrice         string
	Input            string
	Nonce            string
	Value            string
	Date             string
	TransactionIndex string
	Status           string
}

type BigQueryTransaction struct {
	Hash             string `gorm:"unique;not null"` // set to unique and not null
	BlockNumber      string
	BlockHash        string
	Sender           string
	Receiver         string
	Gas              int
	GasPrice         float64
	Input            string
	Nonce            string
	Value            string
	Date             string
	TransactionIndex string
	Status           string
}
