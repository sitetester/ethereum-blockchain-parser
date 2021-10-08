package eth

import (
	"github.com/sitetester/ethereum-blockchain-parser/src/entity/eth"
	"github.com/sitetester/ethereum-blockchain-parser/src/service/eth/client"
	"strconv"
	"strings"
)

const HexToIntDivisor = 1000000000000000000

type BlocksParser struct{}

type TransactionStatus struct {
	Hash   string
	Status string
}

type BlockWithNumber struct {
	Number int
	Block  eth.Block
}

func (blocksParser BlocksParser) ParseBlocks(start int, totalBlock int) []eth.Block {
	var parsedBlocks []eth.Block

	blocksChan := make(chan BlockWithNumber, totalBlock)

	var i int
	for i = start; i < start+totalBlock; i++ {
		go parseBlock(i, blocksChan)
	}

	blocksHashWithStatusMap := getNumberWithBlockMap(blocksChan, totalBlock)

	for _, block := range blocksHashWithStatusMap {
		parsedBlocks = append(parsedBlocks, block)
	}

	return parsedBlocks
}

func parseBlock(blockNumber int, blocksChan chan BlockWithNumber) {
	var infuraClient client.InfuraClient

	parsedBlock := infuraClient.BlockByNumber(blockNumber)
	parsedBlock.EventLogs = infuraClient.GetEventLogs(parsedBlock.Hash)
	adjustBlock(&parsedBlock)

	if len(parsedBlock.Transactions) > 0 {
		setBlockTransactionsStatus(parsedBlock)
	}

	// when block is fully parsed, put it into BlockWithNumber channel
	// Number will be used to identity which block is parsed
	blocksChan <- BlockWithNumber{Number: blockNumber, Block: parsedBlock}
}

func getNumberWithBlockMap(ch chan BlockWithNumber, totalBlocks int) map[int]eth.Block {
	numberWithBlockMap := make(map[int]eth.Block)

	for {
		select {
		case blockWithNumber := <-ch:
			numberWithBlockMap[blockWithNumber.Number] = blockWithNumber.Block

			if len(numberWithBlockMap) == totalBlocks {
				return numberWithBlockMap
			}
		}
	}
}

func setBlockTransactionsStatus(block eth.Block) {
	ch := make(chan TransactionStatus, len(block.Transactions))

	for _, transaction := range block.Transactions {
		go fetchTransactionStatus(transaction.Hash, ch)
	}

	hashWithStatusMap := getHashWithStatusMap(ch, block)

	for i := 0; i < len(block.Transactions); i++ {
		transaction := &block.Transactions[i]
		transaction.Status = hashWithStatusMap[transaction.Hash]
	}
}

func fetchTransactionStatus(hash string, ch chan TransactionStatus) {
	var client client.InfuraClient
	transactionReceipt := client.GetTransactionReceipt(hash)
	if len(transactionReceipt.Status) > 0 {
		ch <- TransactionStatus{Hash: hash, Status: transactionReceipt.Status}
	} else {
		ch <- TransactionStatus{Hash: hash, Status: "not_found"}
	}
}

// https://tour.golang.org/concurrency/5
// https://gobyexample.com/select
// `select` lets us wait/block on channel operations
func getHashWithStatusMap(ch chan TransactionStatus, block eth.Block) map[string]string {
	hashWithStatusMap := make(map[string]string)

	for {
		select {
		case ts := <-ch:
			hashWithStatusMap[ts.Hash] = ts.Status
			if len(hashWithStatusMap) == len(block.Transactions) {
				return hashWithStatusMap
			}
		}
	}
}

func adjustBlock(block *eth.Block) {
	block.TransactionsCount = len(block.Transactions)
	block.NumberInt = hexToInt(block.Number)

	for i := 0; i < len(block.Transactions); i++ {
		transaction := &block.Transactions[i]
		transaction.BlockNumber = hexToIntStr(transaction.BlockNumber)
		transaction.Date = block.Timestamp

		adjustGasValue(*transaction)
		adjustGasPriceValue(*transaction)
		adjustValueInEither(*transaction)
	}

	for i := 0; i < len(block.EventLogs); i++ {
		eventLog := &block.EventLogs[i]
		eventLog.BlockNumber = hexToIntStr(eventLog.BlockNumber)
		eventLog.Date = block.Timestamp
		eventLog.TopicsString = strings.Join(eventLog.Topics, ", ")
	}
}

func hexToInt(hex string) int {
	hexWithout0x := strings.Replace(hex, "0x", "", -1)
	parseInt, err := strconv.ParseInt(hexWithout0x, 16, 64)
	if err != nil {
		panic(err)
	}

	return int(parseInt)
}

func hexToIntStr(s string) string {
	iInt64, _ := strconv.ParseInt(s, 0, 64)
	return strconv.FormatInt(iInt64, 10)
}

func adjustGasValue(transaction eth.Transaction) {
	gasInt, _ := strconv.Atoi(hexToIntStr(transaction.Gas))
	gasPriceStr := (string)(gasInt / HexToIntDivisor)
	transaction.Gas = gasPriceStr
}

func adjustGasPriceValue(transaction eth.Transaction) {
	gasPriceInt, _ := strconv.Atoi(hexToIntStr(transaction.GasPrice))
	gasPriceStr := (string)(gasPriceInt / HexToIntDivisor)
	transaction.GasPrice = gasPriceStr
}

func adjustValueInEither(transaction eth.Transaction) {
	valueInt, _ := strconv.Atoi(hexToIntStr(transaction.Value))
	str := (string)(valueInt / HexToIntDivisor)

	transaction.Value = str
}
