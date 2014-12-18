package tracker

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type connection struct {
	token  string
	client *http.Client
}

func newConnection(token string) connection {
	return connection{
		token:  token,
		client: &http.Client{},
	}
}

func (c connection) Do(request *http.Request, response interface{}) error {
	resp, err := c.sendRequest(request)
	if err != nil {
		return err
	}

	if response != nil {
		return c.decodeResponse(resp, response)
	}

	return nil
}

func (c connection) CreateRequest(method string, path string) (*http.Request, error) {
	request, err := http.NewRequest(method, DefaultURL+"/services/v5"+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}
	request.Header.Add("X-TrackerToken", c.token)

	return request, nil
}

func (c connection) sendRequest(request *http.Request) (*http.Response, error) {
	response, err := c.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %s", err)
	}

	if response.StatusCode == http.StatusUnauthorized {
		return nil, errors.New("invalid token")
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed (%d)", response.StatusCode)
	}

	return response, nil
}

func (c connection) decodeResponse(response *http.Response, object interface{}) error {
	if err := json.NewDecoder(response.Body).Decode(object); err != nil {
		return fmt.Errorf("invalid json response: %s", err)
	}

	err := response.Body.Close()
	if err != nil {
		return fmt.Errorf("error closing response body: %s", err)
	}

	return nil
}
