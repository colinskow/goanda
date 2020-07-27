package goanda

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io/ioutil"
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
	hostname   string
	streamhost string
	token      string
	accountID  string
	headers    map[string]string
	client     *http.Client
}

// ByteStream streams a line of json encoded data or an error
type ByteStream struct {
	Value []byte
	Error error
}

// OandaAgent http agent header
const OandaAgent string = "v20-golang/0.0.1"

// NewConnection makes a new API client
func NewConnection(accountID string, token string, live bool) *OandaConnection {
	var hostname string
	var streamhost string
	// should we use the live API?
	if live {
		hostname = "https://api-fxtrade.oanda.com/v3"
		streamhost = "https://stream-fxtrade.oanda.com/v3"
	} else {
		hostname = "https://api-fxpractice.oanda.com/v3"
		streamhost = "https://stream-fxpractice.oanda.com/v3"
	}
	authHeader := "Bearer " + token

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
		hostname:   hostname,
		streamhost: streamhost,
		token:      token,
		headers:    headers,
		accountID:  accountID,
		client:     &client,
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

// Stream consume a streamng API
func (c *OandaConnection) Stream(endpoint string) (chan ByteStream, context.CancelFunc, error) {
	url := createURL(c.streamhost, endpoint)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	req = req.WithContext(ctx)
	for key, val := range c.headers {
		req.Header.Set(key, val)
	}

	client := http.Client{}
	res, getErr := client.Do(req)
	if getErr != nil {
		cancel()
		return nil, nil, getErr
	}
	reader := bufio.NewReader(res.Body)

	if res.StatusCode != 200 {
		body, readErr := ioutil.ReadAll(res.Body)
		cancel()
		if readErr != nil {
			return nil, nil, readErr
		}
		apiErr := checkAPIErr(body, endpoint)
		if apiErr != nil {
			return nil, nil, apiErr
		}
		return nil, nil, errors.New(res.Status)
	}

	channel := make(chan ByteStream)

	go func() {
		for {
			line, err := reader.ReadBytes('\n')
			if err == nil {
				channel <- ByteStream{Value: line, Error: nil}
			} else {
				channel <- ByteStream{Value: nil, Error: err}
				if ctx.Err() == nil {
					cancel()
				}
				close(channel)
				return
			}
		}
	}()

	return channel, cancel, nil
}
