package goanda

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// CheckErr utility function to log fatal errors
func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkAPIErr(body []byte, route string) error {
	bodyString := string(body[:])
	if strings.Contains(bodyString, "errorMessage") {
		// log.SetFlags(log.LstdFlags | log.Llongfile)
		return errors.New("\nOANDA API Error: " + bodyString + "\nOn route: " + route)
	}
	return nil
}

func unmarshalJSON(body []byte, data interface{}) error {
	jsonErr := json.Unmarshal(body, &data)
	return jsonErr
}

func createURL(host string, endpoint string) string {
	var buffer bytes.Buffer
	// Generate the auth header
	buffer.WriteString(host)
	buffer.WriteString(endpoint)

	url := buffer.String()
	return url
}

func makeRequest(c *OandaConnection, endpoint string, client http.Client, req *http.Request) ([]byte, error) {
	req.Header.Set("User-Agent", c.headers.agent)
	req.Header.Set("Authorization", c.headers.auth)
	req.Header.Set("Content-Type", c.headers.contentType)

	res, getErr := client.Do(req)
	if getErr != nil {
		return nil, getErr
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}
	apiErr := checkAPIErr(body, endpoint)
	if apiErr != nil {
		return nil, apiErr
	}
	return body, nil
}
