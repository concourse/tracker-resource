package tracker

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/xoebus/go-tracker/resources"
)

type ProjectClient struct {
	id   int
	conn connection
}

type State string

const (
	StateUnscheduled = "unscheduled"
	StatePlanned     = "planned"
	StateStarted     = "started"
	StateFinished    = "finished"
	StateDelivered   = "delivered"
	StateAccepted    = "accepted"
	StateRejected    = "rejected"
)

type StoriesQuery struct {
	State State
}

func (query StoriesQuery) Query() url.Values {
	params := url.Values{}
	params.Set("date_format", "millis")

	if query.State != "" {
		params.Set("with_state", string(query.State))
	}

	return params
}

func (p ProjectClient) Stories(query StoriesQuery) (stories []resources.Story, err error) {
	params := query.Query().Encode()
	request, err := p.createRequest("GET", "/stories?"+params)
	if err != nil {
		return stories, err
	}

	err = p.conn.Do(request, &stories)
	return stories, err
}

func (p ProjectClient) DeliverStory(storyId int) error {
	url := fmt.Sprintf("/stories/%d", storyId)
	request, err := p.createRequest("PUT", url)
	if err != nil {
		return err
	}

	p.addJSONBody(request, `{"current_state":"delivered"}`)

	return p.conn.Do(request, nil)
}

func (p ProjectClient) createRequest(method string, path string) (*http.Request, error) {
	projectPath := fmt.Sprintf("/projects/%d%s", p.id, path)
	return p.conn.CreateRequest(method, projectPath)
}

func (p ProjectClient) addJSONBody(request *http.Request, body string) {
	request.Header.Add("Content-Type", "application/json")
	request.Body = ioutil.NopCloser(strings.NewReader(body))
}
