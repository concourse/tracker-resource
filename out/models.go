package out

import (
	"time"
)

type OutRequest struct {
	Source Source `json:"source"`
}

type Source struct {
	Token      string `json:"token"`
	ProjectID  int    `json:"project_id"`
	TrackerURL string `json:"tracker_url"`
}

type OutResponse struct {
	Version Version `json:"version"`
}

type Version struct {
	Time time.Time `json:"time"`
}
