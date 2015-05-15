package in

import "github.com/concourse/tracker-resource"

type InResponse struct {
	Version  resource.Version        `json:"version"`
	Metadata []resource.MetadataPair `json:"metadata"`
}
