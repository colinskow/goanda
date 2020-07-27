package goanda

// Supporting OANDA docs - http://developer.oanda.com/rest-live-v20/instrument-ep/

import (
	"time"

	"github.com/google/go-querystring/query"
)

// Candle https://developer.oanda.com/rest-live-v20/instrument-df/
type Candle struct {
	Open  string `json:"o"`
	Close string `json:"c"`
	Low   string `json:"l"`
	High  string `json:"h"`
}

// Candles https://developer.oanda.com/rest-live-v20/instrument-df/
type Candles struct {
	Complete bool      `json:"complete"`
	Volume   int       `json:"volume"`
	Time     time.Time `json:"time"`
	Mid      Candle    `json:"mid"`
}

// BidAskCandles https://developer.oanda.com/rest-live-v20/instrument-df/
type BidAskCandles struct {
	Candles []struct {
		Ask struct {
			C string `json:"c"`
			H string `json:"h"`
			L string `json:"l"`
			O string `json:"o"`
		} `json:"ask"`
		Bid struct {
			C string `json:"c"`
			H string `json:"h"`
			L string `json:"l"`
			O string `json:"o"`
		} `json:"bid"`
		Complete bool      `json:"complete"`
		Time     time.Time `json:"time"`
		Volume   int       `json:"volume"`
	} `json:"candles"`
}

// InstrumentHistory https://developer.oanda.com/rest-live-v20/instrument-df/
type InstrumentHistory struct {
	Instrument  string    `json:"instrument"`
	Granularity string    `json:"granularity"`
	Candles     []Candles `json:"candles"`
}

// Bucket https://developer.oanda.com/rest-live-v20/instrument-df/
type Bucket struct {
	Price             string `json:"price"`
	LongCountPercent  string `json:"longCountPercent"`
	ShortCountPercent string `json:"shortCountPercent"`
}

// BrokerBook https://developer.oanda.com/rest-live-v20/instrument-df/
type BrokerBook struct {
	Instrument  string    `json:"instrument"`
	Time        time.Time `json:"time"`
	Price       string    `json:"price"`
	BucketWidth string    `json:"bucketWidth"`
	Buckets     []Bucket  `json:"buckets"`
}

// InstrumentPricing https://developer.oanda.com/rest-live-v20/instrument-df/
type InstrumentPricing struct {
	Time   time.Time `json:"time"`
	Prices []struct {
		Type string    `json:"type"`
		Time time.Time `json:"time"`
		Bids []struct {
			Price     string `json:"price"`
			Liquidity int    `json:"liquidity"`
		} `json:"bids"`
		Asks []struct {
			Price     string `json:"price"`
			Liquidity int    `json:"liquidity"`
		} `json:"asks"`
		CloseoutBid    string `json:"closeoutBid"`
		CloseoutAsk    string `json:"closeoutAsk"`
		Status         string `json:"status"`
		Tradeable      bool   `json:"tradeable"`
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
		QuoteHomeConversionFactors struct {
			PositiveUnits string `json:"positiveUnits"`
			NegativeUnits string `json:"negativeUnits"`
		} `json:"quoteHomeConversionFactors"`
		Instrument string `json:"instrument"`
	} `json:"prices"`
}

// GetCandlesPayload https://developer.oanda.com/rest-live-v20/instrument-ep/
type GetCandlesPayload struct {
	Granularity       *string    `url:"granularity,omitempty"`
	Count             *int       `url:"count,omitempty"`
	From              *time.Time `url:"from,omitempty"`
	To                *time.Time `url:"to,omitempty"`
	Smooth            *bool      `url:"smooth,omitempty"`
	IncludeFirst      *bool      `url:"includeFirst,omitempty"`
	DailyAlignment    *int       `url:"dailyAlignment,omitempty"`
	AlignmentTimezone *string    `url:"alignmentTimezone,omitempty"`
	WeeklyAlignment   *string    `url:"weeklyAlignment,omitempty"`
}

// GetCandles https://developer.oanda.com/rest-live-v20/instrument-ep/
func (c *OandaConnection) GetCandles(instrument string, candleSpec GetCandlesPayload) (InstrumentHistory, error) {
	q, _ := query.Values(candleSpec)
	queryStr := q.Encode()
	endpoint := "/instruments/" + instrument + "/candles"
	if len(queryStr) > 0 {
		endpoint = endpoint + "?" + queryStr
	}
	data := InstrumentHistory{}

	candles, err := c.Request(endpoint)
	if err != nil {
		return data, err
	}

	err = unmarshalJSON(candles, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

// GetBidAskCandles https://developer.oanda.com/rest-live-v20/instrument-ep/
func (c *OandaConnection) GetBidAskCandles(instrument string, candleSpec GetCandlesPayload) (BidAskCandles, error) {
	q, _ := query.Values(candleSpec)
	q.Add("price", "ba")
	queryStr := q.Encode()
	endpoint := "/instruments/" + instrument + "/candles" + "?" + queryStr
	data := BidAskCandles{}

	candles, err := c.Request(endpoint)
	if err != nil {
		return data, err
	}

	err = unmarshalJSON(candles, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

// OrderBook https://developer.oanda.com/rest-live-v20/instrument-ep/
func (c *OandaConnection) OrderBook(instrument string) (BrokerBook, error) {
	endpoint := "/instruments/" + instrument + "/orderBook"
	data := BrokerBook{}

	orderbook, err := c.Request(endpoint)
	if err != nil {
		return data, err
	}

	err = unmarshalJSON(orderbook, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

// PositionBook https://developer.oanda.com/rest-live-v20/instrument-ep/
func (c *OandaConnection) PositionBook(instrument string) (BrokerBook, error) {
	endpoint := "/instruments/" + instrument + "/positionBook"
	data := BrokerBook{}
	orderbook, err := c.Request(endpoint)
	if err != nil {
		return data, err
	}

	err = unmarshalJSON(orderbook, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

// GetInstrumentPrice https://developer.oanda.com/rest-live-v20/pricing-ep/
func (c *OandaConnection) GetInstrumentPrice(instruments string) (InstrumentPricing, error) {
	endpoint := "/accounts/" + c.accountID + "/pricing?instruments=" + instruments
	data := InstrumentPricing{}

	pricing, err := c.Request(endpoint)
	if err != nil {
		return data, err
	}

	err = unmarshalJSON(pricing, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}
