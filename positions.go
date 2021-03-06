package goanda

import "encoding/json"

// Supporting OANDA docs - http://developer.oanda.com/rest-live-v20/position-ep/

// OpenPositions https://developer.oanda.com/rest-live-v20/position-df/
type OpenPositions struct {
	LastTransactionID string `json:"lastTransactionID"`
	Positions         []struct {
		Instrument string `json:"instrument"`
		Long       struct {
			AveragePrice string   `json:"averagePrice"`
			Pl           string   `json:"pl"`
			ResettablePL string   `json:"resettablePL"`
			TradeIDs     []string `json:"tradeIDs"`
			Units        string   `json:"units"`
			UnrealizedPL string   `json:"unrealizedPL"`
		} `json:"long"`
		Pl           string `json:"pl"`
		ResettablePL string `json:"resettablePL"`
		Short        struct {
			AveragePrice string   `json:"averagePrice"`
			Pl           string   `json:"pl"`
			ResettablePL string   `json:"resettablePL"`
			TradeIDs     []string `json:"tradeIDs"`
			Units        string   `json:"units"`
			UnrealizedPL string   `json:"unrealizedPL"`
		} `json:"short"`
		UnrealizedPL string `json:"unrealizedPL"`
	} `json:"positions"`
}

// ClosePositionPayload https://developer.oanda.com/rest-live-v20/order-ep/
type ClosePositionPayload struct {
	LongUnits  string `json:"longUnits"`
	ShortUnits string `json:"shortUnits"`
}

// GetOpenPositions https://developer.oanda.com/rest-live-v20/order-ep/
func (c *OandaConnection) GetOpenPositions() (OpenPositions, error) {
	endpoint := "/accounts/" + c.accountID + "/openPositions"
	data := OpenPositions{}

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

// ClosePosition https://developer.oanda.com/rest-live-v20/order-ep/
func (c *OandaConnection) ClosePosition(instrument string, body ClosePositionPayload) (ModifiedTrade, error) {
	endpoint := "/accounts/" + c.accountID + "/positions/" + instrument + "/close"
	data := ModifiedTrade{}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return data, err
	}

	response, err := c.Update(endpoint, jsonBody)
	if err != nil {
		return data, err
	}

	err = unmarshalJSON(response, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}
