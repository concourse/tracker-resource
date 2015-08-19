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

	"github.com/concourse/tracker-resource"
	"github.com/concourse/tracker-resource/out"

	"github.com/mitchellh/colorstring"
	"github.com/xoebus/go-tracker"
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
	if trackerURL == "" {
		trackerURL = "https://www.pivotaltracker.com"
	}

	token := request.Source.Token
	projectID, err := strconv.Atoi(request.Source.ProjectID)
	if err != nil {
		fatal("converting the project ID to an integer", err)
	}

	repos := request.Params.Repos
	sayf("Scanning repositories: %s\n", strings.Join(repos, ", "))

	tracker.DefaultURL = trackerURL
	client := tracker.NewClient(token).InProject(projectID)

	query := tracker.StoriesQuery{
		State: tracker.StoryStateFinished,
	}
	stories, _, err := client.Stories(query)
	if err != nil {
		fatal("getting list of stories", err)
	}

	for _, story := range stories {
		deliverIfDone(client, story, sources, repos)
	}

	outputResponse()
}

func deliverIfDone(client tracker.ProjectClient, story tracker.Story, sources string, repos []string) {
	sayf(colorstring.Color("Checking for finished story: [blue]#%d\n"), story.ID)

	for _, repo := range repos {
		dir := filepath.Join(sources, repo)

		sayf(colorstring.Color("  [white][bold]%s[default]...%s"), repo, strings.Repeat(" ", 80-2-3-10-len(repo)))

		outputFixes := checkGitLog("fixes", story, dir)
		outputFinishes := checkGitLog("finishes", story, dir)

		if len(outputFixes) > 0 || len(outputFinishes) > 0 {
			sayf(colorstring.Color("[green]DELIVERING\n"))
			client.DeliverStoryWithComment(story.ID, "Delivered by Concourse")
		} else {
			sayf(colorstring.Color("  [yellow]SKIPPING\n"))
		}
	}

	sayf("\n")
}

func checkGitLog(verb string, story tracker.Story, dir string) []byte {
	command := exec.Command("git", "log", "-i", "--grep", fmt.Sprintf("%s #%d", verb, story.ID))
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
		Version: resource.Version{
			Time: time.Now(),
		},
	})
}

func sayf(message string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, message, args...)
}

func fatal(doing string, err error) {
	sayf(colorstring.Color("[red]error %s: %s\n"), doing, err)
	os.Exit(1)
}
