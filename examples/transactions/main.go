package main

import (
	"log"
	"os"
	"time"

	"github.com/colinskow/goanda"
	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	key := os.Getenv("OANDA_API_KEY")
	accountID := os.Getenv("OANDA_ACCOUNT_ID")
	oanda := goanda.NewConnection(accountID, key, false)

	toTime := time.Now()
	// AddDate(Y, M, D) - https://golang.org/pkg/time/#Time.AddDate
	fromTime := toTime.AddDate(0, -1, 0)
	transactions, err := oanda.GetTransactions(fromTime, toTime)
	goanda.CheckErr(err)
	spew.Dump("%+v\n", transactions)

	transactionsSince, err := oanda.GetTransactionsSinceID("55")
	goanda.CheckErr(err)
	spew.Dump("%+v\n", transactionsSince)
}
