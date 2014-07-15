package tracker

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/xoebus/go-tracker/resources"
)

var DefaultURL = "https://www.pivotaltracker.com"

type Client struct {
	client *http.Client
	token  string
}

func NewClient(token string) *Client {
	return &Client{
		client: &http.Client{},
		token:  token,
	}
}

func (c Client) Me() (resources.Me, error) {
	var me resources.Me

	request, err := c.createRequest("GET", "/me")
	if err != nil {
		return me, err
	}

	response, err := c.sendRequest(request)
	if err != nil {
		return me, err
	}

	if err := c.decodeResponse(response, &me); err != nil {
		return me, err
	}

	return me, nil
}

func (c Client) InProject(projectId int) ProjectClient {
	return ProjectClient{
		id:     projectId,
		client: c,
	}
}

func (c Client) createRequest(method string, path string) (*http.Request, error) {
	request, err := http.NewRequest(method, DefaultURL+"/services/v5"+path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}
	request.Header.Add("X-TrackerToken", c.token)

	return request, nil
}

func (c Client) sendRequest(request *http.Request) (*http.Response, error) {
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

func (c Client) decodeResponse(response *http.Response, object interface{}) error {
	if err := json.NewDecoder(response.Body).Decode(object); err != nil {
		return fmt.Errorf("invalid json response: %s", err)
	}

	err := response.Body.Close()
	if err != nil {
		return fmt.Errorf("error closing response body: %s", err)
	}

	return nil
}
