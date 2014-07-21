package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/concourse/tracker-resource/out"

	"github.com/xoebus/go-tracker"
	"github.com/xoebus/go-tracker/resources"
)

func buildRequest() out.OutRequest {
	var request out.OutRequest

	err := json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		fatal("reading request", err)
	}

	return request
}

func main() {
	if len(os.Args) < 2 {
		sayf("usage: %s <sources directory>\n", os.Args[0])
		os.Exit(1)
	}

	sources := os.Args[1]
	request := buildRequest()

	trackerURL := request.Source.TrackerURL
	token := request.Source.Token
	projectID, err := strconv.Atoi(request.Source.ProjectID)
	if err != nil {
		fatal("converting the project ID to an integer", err)
	}

	repos := request.Params.Repos
	sayf("Scanning repositories: %s\n", strings.Join(repos, ", "))

	tracker.DefaultURL = trackerURL
	client := tracker.NewClient(token).InProject(projectID)
	stories, err := client.Stories()
	if err != nil {
		fatal("getting list of stories", err)
	}

	for _, story := range stories {
		deliverIfDone(client, story, sources, repos)
	}

	outputResponse()
}

func deliverIfDone(client tracker.ProjectClient, story resources.Story, sources string, repos []string) {
	for _, repo := range repos {
		dir := filepath.Join(sources, repo)

		sayf("checking in dir %s\n", dir)

		outputFixes := checkOutput("fixes", story, dir)
		outputFinishes := checkOutput("finishes", story, dir)

		if len(outputFixes) > 0 || len(outputFinishes) > 0 {
			sayf("found the story, delivering it!: %d\n", story.ID)
			client.DeliverStory(story.ID)
		} else {
			sayf("could not find story for delivery: %d\n", story.ID)
		}
	}
}

func checkOutput(verb string, story resources.Story, dir string) []byte {
	command := exec.Command("git", "log", "--grep", fmt.Sprintf("%s #%d", verb, story.ID))
	command.Dir = dir

	output, err := command.CombinedOutput()
	if err != nil {
		sayf("git logging failed for story: %d: %s\n", story.ID, err)

		return nil
	}

	return output
}

func outputResponse() {
	json.NewEncoder(os.Stdout).Encode(out.OutResponse{
		Version: out.Version{
			Time: time.Now(),
		},
	})
}

func sayf(message string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, message, args...)
}

func fatal(doing string, err error) {
	sayf("error %s: %s\n", doing, err)
	os.Exit(1)
}
