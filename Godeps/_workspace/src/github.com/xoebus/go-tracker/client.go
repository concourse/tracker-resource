package tracker

import "github.com/xoebus/go-tracker/resources"

var DefaultURL = "https://www.pivotaltracker.com"

type Client struct {
	conn connection
}

func NewClient(token string) *Client {
	return &Client{
		conn: newConnection(token),
	}
}

func (c Client) Me() (me resources.Me, err error) {
	request, err := c.conn.CreateRequest("GET", "/me")
	if err != nil {
		return me, err
	}

	err = c.conn.Do(request, &me)

	return me, err
}

func (c Client) InProject(projectId int) ProjectClient {
	return ProjectClient{
		id:   projectId,
		conn: c.conn,
	}
}
