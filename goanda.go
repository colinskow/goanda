package goanda

import (
	"bytes"
	"net/http"
	"time"
)

// Headers data
type Headers struct {
	contentType    string
	agent          string
	DatetimeFormat string
	auth           string
}

// Connection methods
type Connection interface {
	Request(endpoint string) []byte
	Send(endpoint string, data []byte) []byte
	Update(endpoint string, data []byte) []byte
	GetOrderDetails(instrument string, units string) OrderDetails
	GetAccountSummary() AccountSummary
	CreateOrder(body OrderPayload) OrderResponse
}

// OandaConnection data
type OandaConnection struct {
	hostname       string
	port           int
	ssl            bool
	token          string
	accountID      string
	DatetimeFormat string
	headers        map[string]string
	client         *http.Client
}

// OandaAgent http agent header
const OandaAgent string = "v20-golang/0.0.1"

// NewConnection makes a new API client
func NewConnection(accountID string, token string, live bool) *OandaConnection {
	hostname := ""
	// should we use the live API?
	if live {
		hostname = "https://api-fxtrade.oanda.com/v3"
	} else {
		hostname = "https://api-fxpractice.oanda.com/v3"
	}

	var buffer bytes.Buffer
	// Generate the auth header
	buffer.WriteString("Bearer ")
	buffer.WriteString(token)

	authHeader := buffer.String()
	// Create headers for oanda to be used in requests
	headers := map[string]string{
		"User-Agent":    OandaAgent,
		"Authorization": authHeader,
		"Content-Type":  "application/json",
		"Connection":    "Keep-Alive",
	}
	client := http.Client{
		Timeout: time.Second * 5, // 5 sec timeout
	}
	// Create the connection object
	connection := &OandaConnection{
		hostname:  hostname,
		port:      443,
		ssl:       true,
		token:     token,
		headers:   headers,
		accountID: accountID,
		client:    &client,
	}

	return connection
}

// Request make a get request to the Oanda V20 API
// TODO: include params as a second option
func (c *OandaConnection) Request(endpoint string) ([]byte, error) {

	url := createURL(c.hostname, endpoint)

	// New request object
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	body, err := makeRequest(c, endpoint, c.client, req)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Send post data to the API
func (c *OandaConnection) Send(endpoint string, data []byte) ([]byte, error) {
	url := createURL(c.hostname, endpoint)

	// New request object
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	body, err := makeRequest(c, endpoint, c.client, req)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Update put data to the API
func (c *OandaConnection) Update(endpoint string, data []byte) ([]byte, error) {
	url := createURL(c.hostname, endpoint)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	body, err := makeRequest(c, endpoint, c.client, req)
	if err != nil {
		return nil, err
	}

	return body, nil
}
