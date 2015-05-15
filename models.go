package resource

import "time"

type Source struct {
	Token      string `json:"token"`
	ProjectID  string `json:"project_id"`
	TrackerURL string `json:"tracker_url"`
}

type Version struct {
	Time time.Time `json:"time"`
}

type MetadataPair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
