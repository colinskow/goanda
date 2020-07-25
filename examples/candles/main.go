package main

import (
	"fmt"
	"log"
	"os"

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
	fmt.Println(key, accountID)
	oanda := goanda.NewConnection(accountID, key, false)
	gran := "M5"
	count := 60
	candleSpec := goanda.GetCandlesPayload{
		Granularity: &gran,
		Count:       &count}
	history, err := oanda.GetCandles("EUR_USD", candleSpec)
	goanda.CheckErr(err)
	spew.Dump(history)
}
