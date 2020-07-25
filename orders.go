package goanda

// Supporting OANDA docs - http://developer.oanda.com/rest-live-v20/order-ep/

import (
	"encoding/json"
	"time"
)

// OrderExtensions https://developer.oanda.com/rest-live-v20/order-df/
type OrderExtensions struct {
	Comment string `json:"comment,omitempty"`
	ID      string `json:"id,omitempty"`
	Tag     string `json:"tag,omitempty"`
}

// OnFill https://developer.oanda.com/rest-live-v20/order-df/
type OnFill struct {
	TimeInForce string `json:"timeInForce,omitempty"`
	Price       string `json:"price,omitempty"` // must be a string for float precision
}

// OrderBody https://developer.oanda.com/rest-live-v20/order-df/
type OrderBody struct {
	Units            int              `json:"units"`
	Instrument       string           `json:"instrument"`
	TimeInForce      string           `json:"timeInForce"`
	Type             string           `json:"type"`
	PositionFill     string           `json:"positionFill,omitempty"`
	Price            string           `json:"price,omitempty"`
	TakeProfitOnFill *OnFill          `json:"takeProfitOnFill,omitempty"`
	StopLossOnFill   *OnFill          `json:"stopLossOnFill,omitempty"`
	ClientExtensions *OrderExtensions `json:"clientExtensions,omitempty"`
	TradeID          string           `json:"tradeId,omitempty"`
}

// OrderPayload https://developer.oanda.com/rest-live-v20/order-df/
type OrderPayload struct {
	Order OrderBody `json:"order"`
}

// OrderResponse https://developer.oanda.com/rest-live-v20/order-df/
type OrderResponse struct {
	LastTransactionID      string `json:"lastTransactionID"`
	OrderCreateTransaction struct {
		AccountID    string    `json:"accountID"`
		BatchID      string    `json:"batchID"`
		ID           string    `json:"id"`
		Instrument   string    `json:"instrument"`
		PositionFill string    `json:"positionFill"`
		Reason       string    `json:"reason"`
		Time         time.Time `json:"time"`
		TimeInForce  string    `json:"timeInForce"`
		Type         string    `json:"type"`
		Units        string    `json:"units"`
		UserID       int       `json:"userID"`
	} `json:"orderCreateTransaction"`
	OrderFillTransaction struct {
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
	} `json:"orderFillTransaction"`
	RelatedTransactionIDs []string `json:"relatedTransactionIDs"`
}

// OrderInfo https://developer.oanda.com/rest-live-v20/order-df/
type OrderInfo struct {
	ClientExtensions struct {
		Comment string `json:"comment,omitempty"`
		ID      string `json:"id,omitempty"`
		Tag     string `json:"tag,omitempty"`
	} `json:"clientExtensions,omitempty"`
	CreateTime       time.Time `json:"createTime"`
	ID               string    `json:"id"`
	Instrument       string    `json:"instrument,omitempty"`
	PartialFill      string    `json:"partialFill"`
	PositionFill     string    `json:"positionFill"`
	Price            string    `json:"price"`
	ReplacesOrderID  string    `json:"replacesOrderID,omitempty"`
	State            string    `json:"state"`
	TimeInForce      string    `json:"timeInForce"`
	TriggerCondition string    `json:"triggerCondition"`
	Type             string    `json:"type"`
	Units            string    `json:"units,omitempty"`
}

// RetrievedOrders https://developer.oanda.com/rest-live-v20/order-ep/
type RetrievedOrders struct {
	LastTransactionID string      `json:"lastTransactionID"`
	Orders            []OrderInfo `json:"orders,omitempty"`
}

// RetrievedOrder https://developer.oanda.com/rest-live-v20/order-ep/
type RetrievedOrder struct {
	Order OrderInfo `json:"order"`
}

// CancelledOrder https://developer.oanda.com/rest-live-v20/order-df/
type CancelledOrder struct {
	OrderCancelTransaction struct {
		ID                string    `json:"id"`
		Time              time.Time `json:"time"`
		UserID            int       `json:"userID"`
		AccountID         string    `json:"accountID"`
		BatchID           string    `json:"batchID"`
		RequestID         string    `json:"requestID"`
		Type              string    `json:"type"`
		OrderID           string    `json:"orderID"`
		ClientOrderID     string    `json:"clientOrderID"`
		Reason            string    `json:"reason"`
		ReplacedByOrderID string    `json:"replacedByOrderID"`
	} `json:"orderCancelTransaction"`
	RelatedTransactionIDs []string `json:"relatedTransactionIDs"`
	LastTransactionID     string   `json:"lastTransactionID"`
}

// CreateOrder https://developer.oanda.com/rest-live-v20/order-ep/
func (c *OandaConnection) CreateOrder(body OrderPayload) (OrderResponse, error) {
	endpoint := "/accounts/" + c.accountID + "/orders"
	data := OrderResponse{}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return data, err
	}

	response, err := c.Send(endpoint, jsonBody)
	if err != nil {
		return data, err
	}

	err = unmarshalJSON(response, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

// GetOrders https://developer.oanda.com/rest-live-v20/order-ep/
func (c *OandaConnection) GetOrders(instrument string) (RetrievedOrders, error) {
	endpoint := "/accounts/" + c.accountID + "/orders"

	if instrument != "" {
		endpoint = endpoint + "?instrument=" + instrument
	}

	data := RetrievedOrders{}
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

// GetPendingOrders https://developer.oanda.com/rest-live-v20/order-ep/
func (c *OandaConnection) GetPendingOrders() (RetrievedOrders, error) {
	endpoint := "/accounts/" + c.accountID + "/pendingOrders"
	data := RetrievedOrders{}

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

// GetOrder https://developer.oanda.com/rest-live-v20/order-ep/
func (c *OandaConnection) GetOrder(orderSpecifier string) (RetrievedOrder, error) {
	endpoint := "/accounts/" + c.accountID + "/orders/" + orderSpecifier
	data := RetrievedOrder{}

	response, err := c.Request(endpoint)
	if err != nil {
		return data, err
	}

	err = unmarshalJSON(response, &data)
	if err != nil {
		return data, err
	}

	return data, err
}

// UpdateOrder https://developer.oanda.com/rest-live-v20/order-ep/
func (c *OandaConnection) UpdateOrder(orderSpecifier string, body OrderPayload) (RetrievedOrder, error) {
	endpoint := "/accounts/" + c.accountID + "/orders/" + orderSpecifier
	data := RetrievedOrder{}

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

// CancelOrder https://developer.oanda.com/rest-live-v20/order-ep/
func (c *OandaConnection) CancelOrder(orderSpecifier string) (CancelledOrder, error) {
	endpoint := "/accounts/" + c.accountID + "/orders/" + orderSpecifier + "/cancel"
	data := CancelledOrder{}

	response, err := c.Update(endpoint, nil)
	if err != nil {
		return data, err
	}

	err = unmarshalJSON(response, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}
