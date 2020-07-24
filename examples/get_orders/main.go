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
	orders, err := oanda.GetOrders("EUR_USD")
	goanda.CheckErr(err)
	pendingOrders, err := oanda.GetPendingOrders()
	goanda.CheckErr(err)
	spew.Dump("%+v\n", orders)
	spew.Dump("%+v\n", pendingOrders)
}
