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
	order := goanda.OrderPayload{
		Order: goanda.OrderBody{
			Units:        100,
			Instrument:   "EUR_USD",
			TimeInForce:  "FOK",
			Type:         "LIMIT",
			PositionFill: "DEFAULT",
			Price:        "1.25000",
			StopLossOnFill: &goanda.OnFill{
				TimeInForce: "GTC",
				Price:       "1.24000",
			},
		},
	}
	orderResult, err := oanda.CreateOrder(order)
	goanda.CheckErr(err)
	spew.Dump("%+v\n", orderResult)
}
