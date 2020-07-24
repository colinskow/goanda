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

	accountChanges, err := oanda.GetAccountChanges("101-011-6559702-001", "54")
	goanda.CheckErr(err)
	spew.Dump(accountChanges)

	accountInstruments, err := oanda.GetAccountInstruments("101-011-6559702-001")
	goanda.CheckErr(err)
	spew.Dump(accountInstruments)

	accountSummary, err := oanda.GetAccountSummary()
	goanda.CheckErr(err)
	spew.Dump(accountSummary)

	account, err := oanda.GetAccount("101-011-6559702-003")
	goanda.CheckErr(err)
	spew.Dump(account)

	accounts, err := oanda.GetAccounts()
	goanda.CheckErr(err)
	spew.Dump(accounts)
}
