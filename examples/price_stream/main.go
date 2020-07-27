package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/colinskow/goanda"
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

	instruments := []string{"EUR_USD", "GBP_USD"}
	priceStream, cancel, err := oanda.StreamPricing(instruments)
	goanda.CheckErr(err)

	// Quit after 30 seconds of listening
	timer := time.NewTimer(30 * time.Second)
	go func() {
		<-timer.C
		cancel()
	}()

	for item := range priceStream {
		if item.Error == nil {
			fmt.Println(item.ClientPrice.Instrument, item.ClientPrice.Bids[0].Price, item.ClientPrice.Asks[0].Price)
		} else {
			fmt.Println(item.Error)
		}
	}
}
