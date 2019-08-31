package eth

import (
	"cloud.google.com/go/bigquery"
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	"log"
	"os"
)

const ProjectID = "insert-biquery-project-id-here"
const Dataset = "blockchains"

const TransactionsTable = "eth_transactions"

type BigQueryClient struct {
}

func (bqClient BigQueryClient) GetClient() *bigquery.Client {
	proj := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if proj == "" {
		fmt.Println("GOOGLE_APPLICATION_CREDENTIALS environment variable must be set.")
		os.Exit(1)
	}

	ctx := context.Background()
	client, err := bigquery.NewClient(ctx, ProjectID)
	if err != nil {
		log.Fatal(err)
	}

	return client
}

func (bqClient BigQueryClient) RunQuery(query string) {
	q := bqClient.GetClient().Query("SELECT `blockNumber` from `blockchains.dev_timestamp_partitioned_eth_event_logs` WHERE date(date) = '2019-08-26' ORDER BY blockNumber DESC LIMIT 1")

	ctx := context.Background()
	it, err := q.Read(ctx)
	if err != nil {
		log.Fatal(err) // TODO: Handle error.
	}

	for {
		var values []bigquery.Value
		err := it.Next(&values)
		if err == iterator.Done {
			break
		}
		if err != nil {
			// TODO: Handle error.
		}
		fmt.Println(values)
	}

}

func (bqClient BigQueryClient) InsertRows(tableId string, src interface{}) bool {
	myDataset := bqClient.GetClient().Dataset(Dataset)
	table := myDataset.Table(tableId)
	inserter := table.Inserter()

	// Item implements the ValueSaver interface.
	/*EthEventLogData := []*EventLogData{
		{Address: "b4", BlockHash: "b4", BlockNumber: "444", Date: (int)(time.Now().Unix()), Transactionlogindex: "bbbbTransactionlogindex444", Type: "babc4"},
	}*/

	ctx := context.Background()
	err := inserter.Put(ctx, src)
	if err != nil {
		if multiError, ok := err.(bigquery.PutMultiError); ok {
			for _, err1 := range multiError {
				for _, err2 := range err1.Errors {
					fmt.Println(err2)
				}
			}
		} else {
			fmt.Println(err)
		}

		return false
	}

	return true
}

type EventLogData struct {
	Address             string
	BlockHash           string
	BlockNumber         string
	Data                string
	LogIndex            string
	Removed             bool
	TransactionHash     string
	TransactionIndex    string
	Date                int
	Transactionlogindex string
	Type                string
}
