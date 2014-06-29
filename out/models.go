package out

import (
	"time"
)

type OutRequest struct {
	Params Params `json:"params"`
}

type Params struct {
	Token     string `json:"token"`
	ProjectID int    `json:"project_id"`
}

type OutResponse struct {
	Version Version `json:"version"`
}

type Version struct {
	Time time.Time `json:"time"`
}