package goanda

import (
	"bytes"
	"context"
	"net/url"
	"strings"
	"time"
)

// Supporting OANDA docs - http://developer.oanda.com/rest-live-v20/pricing-ep/

// ClientPrice https://developer.oanda.com/rest-live-v20/pricing-df/
type ClientPrice struct {
	Asks []struct {
		Liquidity int    `json:"liquidity"`
		Price     string `json:"price"`
	} `json:"asks"`
	Bids []struct {
		Liquidity int    `json:"liquidity"`
		Price     string `json:"price"`
	} `json:"bids"`
	CloseoutAsk                string `json:"closeoutAsk"`
	CloseoutBid                string `json:"closeoutBid"`
	Instrument                 string `json:"instrument"`
	QuoteHomeConversionFactors struct {
		NegativeUnits string `json:"negativeUnits"`
		PositiveUnits string `json:"positiveUnits"`
	} `json:"quoteHomeConversionFactors"`
	Status         string    `json:"status"`
	Time           time.Time `json:"time"`
	UnitsAvailable struct {
		Default struct {
			Long  string `json:"long"`
			Short string `json:"short"`
		} `json:"default"`
		OpenOnly struct {
			Long  string `json:"long"`
			Short string `json:"short"`
		} `json:"openOnly"`
		ReduceFirst struct {
			Long  string `json:"long"`
			Short string `json:"short"`
		} `json:"reduceFirst"`
		ReduceOnly struct {
			Long  string `json:"long"`
			Short string `json:"short"`
		} `json:"reduceOnly"`
	} `json:"unitsAvailable"`
}

// Pricings https://developer.oanda.com/rest-live-v20/pricing-df/
type Pricings struct {
	Prices []ClientPrice `json:"prices"`
}

// PricingStream used to stream prices over a go channel
type PricingStream struct {
	ClientPrice *ClientPrice
	Error       error
}

// GetPricingForInstruments https://developer.oanda.com/rest-live-v20/pricing-ep/
func (c *OandaConnection) GetPricingForInstruments(instruments []string) (Pricings, error) {
	instrumentString := strings.Join(instruments, ",")
	endpoint := "/accounts/" + c.accountID + "/pricing?instruments=" + url.QueryEscape(instrumentString)
	data := Pricings{}

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

// StreamPricing https://developer.oanda.com/rest-live-v20/pricing-ep/
func (c *OandaConnection) StreamPricing(instruments []string) (chan PricingStream, context.CancelFunc, error) {
	instrumentString := strings.Join(instruments, ",")
	endpoint := "/accounts/" + c.accountID + "/pricing/stream?instruments=" + url.QueryEscape(instrumentString)
	output := make(chan PricingStream)

	stream, cancel, err := c.Stream(endpoint)
	if err != nil {
		return nil, nil, err
	}

	go func() {
		for item := range stream {
			if item.Error == nil {
				if bytes.Contains(item.Value, []byte("PRICE")) {
					data := ClientPrice{}
					err = unmarshalJSON(item.Value, &data)
					if err == nil {
						output <- PricingStream{ClientPrice: &data, Error: nil}
					} else {
						output <- PricingStream{ClientPrice: nil, Error: err}
					}
				}
			} else {
				output <- PricingStream{nil, item.Error}
				close(output)
				break
			}
		}
	}()

	return output, cancel, nil
}
