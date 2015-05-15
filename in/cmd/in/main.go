package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/concourse/tracker-resource"
	"github.com/concourse/tracker-resource/in"
)

func main() {
	err := json.NewEncoder(os.Stdout).Encode(
		in.InResponse{
			Version: resource.Version{
				Time: time.Now(),
			},
		},
	)

	if err != nil {
		fmt.Fprintln(os.Stderr, "error writing response: ", err)
		os.Exit(1)
	}
}
