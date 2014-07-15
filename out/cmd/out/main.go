package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/concourse/tracker-resource/out"

	"github.com/xoebus/go-tracker"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <sources directory>\n", os.Args[0])
		os.Exit(1)
	}

	// sources := os.Args[1]

	var request out.OutRequest

	err := json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		fatal("reading request", err)
	}

	trackerURL := request.Params.TrackerURL
	fmt.Fprintf(os.Stderr, "Tracker URL: %s\n", trackerURL)

	token := request.Params.Token
	fmt.Fprintf(os.Stderr, "Tracker Token: %s\n", token)

	projectID := request.Params.ProjectID
	fmt.Fprintf(os.Stderr, "Tracker Project ID: %d\n", projectID)

	tracker.DefaultURL = trackerURL
	client := tracker.NewClient(token)

	stories, err := client.InProject(projectID).Stories()
	if err != nil {
		fatal("getting list of stories", err)
	}

	fmt.Fprintf(os.Stderr, "%+v", stories)

	outputResponse()
}

func outputResponse() {
	json.NewEncoder(os.Stdout).Encode(out.OutResponse{
		Version: out.Version{
			Time: time.Now(),
		},
	})
}

func fatal(doing string, err error) {
	fmt.Fprintf(os.Stderr, "error %s: %s\n", doing, err)
	os.Exit(1)
}
