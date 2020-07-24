package main

import (
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
	oanda := goanda.NewConnection(accountID, key, false)

	instruments := []string{"AUD_USD", "EUR_NZD"}
	orderResponse, err := oanda.GetPricingForInstruments(instruments)
	goanda.CheckErr(err)
	spew.Dump("%+v\n", orderResponse)
}
