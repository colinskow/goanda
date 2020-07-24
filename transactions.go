package goanda

import (
	"net/url"
	"time"
)

// Supporting OANDA docs - http://developer.oanda.com/rest-live-v20/transaction-ep/

// TransactionPages http://developer.oanda.com/rest-live-v20/transaction-df/
type TransactionPages struct {
	Count             int       `json:"count"`
	From              time.Time `json:"from"`
	LastTransactionID string    `json:"lastTransactionID"`
	PageSize          int       `json:"pageSize"`
	Pages             []string  `json:"pages"`
	To                time.Time `json:"to"`
}

// Transaction http://developer.oanda.com/rest-live-v20/transaction-df/
type Transaction struct {
	LastTransactionID string `json:"lastTransactionID"`
	Transaction       struct {
		AccountBalance string    `json:"accountBalance"`
		AccountID      string    `json:"accountID"`
		BatchID        string    `json:"batchID"`
		Financing      string    `json:"financing"`
		ID             string    `json:"id"`
		Instrument     string    `json:"instrument"`
		OrderID        string    `json:"orderID"`
		Pl             string    `json:"pl"`
		Price          string    `json:"price"`
		Reason         string    `json:"reason"`
		Time           time.Time `json:"time"`
		TradeOpened    struct {
			TradeID string `json:"tradeID"`
			Units   string `json:"units"`
		} `json:"tradeOpened"`
		Type   string `json:"type"`
		Units  string `json:"units"`
		UserID int    `json:"userID"`
	} `json:"transaction"`
}

// Transactions http://developer.oanda.com/rest-live-v20/transaction-df/
type Transactions struct {
	LastTransactionID string `json:"lastTransactionID"`
	Transactions      []struct {
		AccountBalance string    `json:"accountBalance"`
		AccountID      string    `json:"accountID"`
		BatchID        string    `json:"batchID"`
		Financing      string    `json:"financing"`
		ID             string    `json:"id"`
		Instrument     string    `json:"instrument"`
		OrderID        string    `json:"orderID"`
		Pl             string    `json:"pl"`
		Price          string    `json:"price"`
		Reason         string    `json:"reason"`
		Time           time.Time `json:"time"`
		TradeOpened    struct {
			TradeID string `json:"tradeID"`
			Units   string `json:"units"`
		} `json:"tradeOpened"`
		Type   string `json:"type"`
		Units  string `json:"units"`
		UserID int    `json:"userID"`
	} `json:"transactions"`
}

// GetTransactions http://developer.oanda.com/rest-live-v20/transaction-ep/
// https://golang.org/pkg/time/#Time.AddDate
// https://play.golang.org/p/Dw7D4JJ7EC
func (c *OandaConnection) GetTransactions(from time.Time, to time.Time) (TransactionPages, error) {
	toTime := to.Format(time.RFC3339)
	fromTime := from.Format(time.RFC3339)
	endpoint := "/accounts/" + c.accountID + "/transactions?to=" + url.QueryEscape(toTime) + "&from=" + url.QueryEscape(fromTime)
	data := TransactionPages{}

	response, err := c.Request(endpoint)
	if err != nil {
		return data, err
	}

	err = unmarshalJSON(response, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}

// GetTransaction http://developer.oanda.com/rest-live-v20/transaction-ep/
func (c *OandaConnection) GetTransaction(ticket string) (Transaction, error) {
	endpoint := "/accounts/" + c.accountID + "/transactions/" + ticket
	data := Transaction{}

	response, err := c.Request(endpoint)
	if err != nil {
		return data, err
	}

	err = unmarshalJSON(response, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

// GetTransactionsSinceID http://developer.oanda.com/rest-live-v20/transaction-ep/
func (c *OandaConnection) GetTransactionsSinceID(id string) (Transactions, error) {

	endpoint := "/accounts/" + c.accountID + "/transactions/sinceid?id=" + id
	data := Transactions{}

	response, err := c.Request(endpoint)
	if err != nil {
		return data, err
	}

	err = unmarshalJSON(response, &data)
	if err != nil {
		return data, err
	}
	return data, nil
}
