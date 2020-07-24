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

	tradesResponse, err := oanda.GetTrade("54")
	goanda.CheckErr(err)
	spew.Dump("%+v\n", tradesResponse)

	openTrades, err := oanda.GetOpenTrades()
	goanda.CheckErr(err)
	spew.Dump("%+v\n", openTrades)

	trade, err := oanda.GetTradesForInstrument("AUD_USD")
	goanda.CheckErr(err)
	spew.Dump("%+v\n", trade)

	reduceTrade, err := oanda.ReduceTradeSize("AUD_USD", goanda.CloseTradePayload{
		Units: "100.00",
	})
	goanda.CheckErr(err)

	spew.Dump("%+v\n", reduceTrade)
}
