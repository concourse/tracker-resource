package out

import (
	"time"
)

type OutRequest struct {
	Source Source `json:"source"`
	Params Params `json:"params"`
}

type Params struct {
	Repos []string `json:"repos"`
}

type Source struct {
	Token      string `json:"token"`
	ProjectID  string `json:"project_id"`
	TrackerURL string `json:"tracker_url"`
}

type OutResponse struct {
	Version Version `json:"version"`
}

type Version struct {
	Time time.Time `json:"time"`
}
