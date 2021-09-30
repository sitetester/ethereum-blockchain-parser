package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"github.com/sitetester/ethereum-blockchain-parser/src/entity/eth"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type InfuraClient struct {
}

const URL = "https://mainnet.infura.io/v3/c8d36b72d2d04f16a94931809cdf6383"

type LatestBlockResponse struct {
}

type RpcResponse struct {
	Jsonrpc string
	Id      int
	Result  string
}

func (client InfuraClient) LatestBlock() int64 {
	var data = []byte(`{"jsonrpc":"2.0","method":"eth_blockNumber","params": [],"id":1}`)
	var rpcResponse RpcResponse
	json.Unmarshal(makeRequest(data), &rpcResponse)
	d, _ := strconv.ParseInt(rpcResponse.Result, 0, 64)

	return d
}

func (client InfuraClient) LatestBlockNumber() string {
	var jsonStr = []byte(`{"jsonrpc":"2.0","method":"eth_blockNumber","params": [],"id":1}`)
	var rpcResponse RpcResponse
	json.Unmarshal(makeRequest(jsonStr), &rpcResponse)

	return rpcResponse.Result
}

type RpcParams struct {
	Jsonrpc string
	Method  string
	Params  []interface{}
	Id      int
}

type BlockResponse struct {
	Jsonrpc string
	Id      int
	Result  eth.Block
}

type BlockByNumberRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []interface{}
	Id      int
}

// https://socketloop.com/tutorials/golang-convert-integer-to-binary-octal-hexadecimal-and-back-to-integer
func hex(i int) string {
	i64 := int64(i)
	str := strings.ToUpper(strconv.FormatInt(i64, 16))
	return "0x" + str
}

// BlockByNumber https://etherscan.io/block/3415124 - block(3415136) with 2 transactions
// https://etherscan.io/block/3415124 - block(3415124) with 5 transactions
// -d '{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params": ["0x5BAD55",false],"id":1}'
func (client InfuraClient) BlockByNumber(blockNumber int) eth.Block {
	params := []interface{}{hex(blockNumber), true}
	b := &BlockByNumberRequest{"2.0", "eth_getBlockByNumber", params, 1}
	out, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}

	jsonStr := string(out)
	// fmt.Println(jsonStr)
	var blockResponse BlockResponse
	if err := json.Unmarshal(makeRequest([]byte(jsonStr)), &blockResponse); err != nil {
		panic(err)
	}

	// fmt.Printf("%+v", blockResponse.Result)

	return blockResponse.Result
}

type TransactionReceiptRequest struct {
	Jsonrpc string
	Method  string
	Params  []string
	Id      int
}

type TransactionReceiptResponse struct {
	Jsonrpc string
	Method  string
	Result  TransactionReceipt
}

type TransactionReceipt struct {
	transactionHash string
	Status          string
}

// https://infura.io/docs/ethereum/json-rpc/eth_getTransactionByHash
// -d '{"jsonrpc":"2.0","method":"eth_getTransactionByHash","params": ["0xbb3a336e3f823ec18197f1e13ee875700f08f03e2cab75f0d0b118dabb44cba0"],"id":1}'
func (client InfuraClient) GetTransactionReceipt(hash string) TransactionReceipt {
	// hash = "0xbb3a336e3f823ec18197f1e13ee875700f08f03e2cab75f0d0b118dabb44cba0"
	params := []string{hash}
	b := &TransactionReceiptRequest{"2.0", "eth_getTransactionReceipt", params, 1}
	out, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}

	jsonStr := string(out)
	// fmt.Println(jsonStr)

	var r TransactionReceiptResponse
	if err := json.Unmarshal(makeRequest([]byte(jsonStr)), &r); err != nil {
		panic(err)
	}

	return r.Result
}

type GetLogsRequest struct {
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  []interface{}
	Id      int
}

type EventLogsResponse struct {
	Jsonrpc string
	Method  string
	Result  []eth.EventLog
}

func (client InfuraClient) GetEventLogs(blockHash string) []eth.EventLog {
	/*params := []interface{}{
		map[string]string{"blockHash": "0xfc66cc2e39c1537a84e290225f2046dfff565ad6fffe36bda2bb24593e0b6a02"},
	}*/

	params := []interface{}{
		map[string]string{"blockHash": blockHash},
	}

	b := &GetLogsRequest{"2.0", "eth_getLogs", params, 1}
	out, err := json.Marshal(b)
	if err != nil {
		panic(err)
	}

	jsonStr := string(out)
	// fmt.Println(jsonStr)

	var eventLogsResponse EventLogsResponse
	// var eventLogsResponse interface{}
	if err := json.Unmarshal(makeRequest([]byte(jsonStr)), &eventLogsResponse); err != nil {
		panic(err)
	}

	// fmt.Println(eventLogsResponse)
	// os.Exit(1)
	return eventLogsResponse.Result
}

func makeRequest(data []byte) []byte {
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Create New http Transport
	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // disable verify
	}

	httpClient := &http.Client{
		Timeout:   time.Duration(1 * time.Hour),
		Transport: transCfg,
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return body
}
