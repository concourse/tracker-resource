package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	commentPath := request.Params.CommentPath

	var comment []byte
	if commentPath != "" {
		if comment, err = ioutil.ReadFile(filepath.Join(sources, commentPath)); err != nil {
			fatal("reading comment file", err)
		}
	}

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
		deliverIfDone(client, story, sources, string(comment), repos)
	}

	outputResponse()
}

func deliverIfDone(client tracker.ProjectClient, story tracker.Story, sources, comment string, repos []string) {
	sayf(colorstring.Color("Checking for finished story: [blue]#%d\n"), story.ID)

	storyResetTime, err := getStoryResetTime(client, story)
	if err != nil {
		fatal(fmt.Sprintf("fetching activity for story #%d", story.ID), err)
	}

	for _, repo := range repos {
		dir := filepath.Join(sources, repo)

		sayf(colorstring.Color("  [white][bold]%s[default]...%s"), repo, strings.Repeat(" ", 80-2-3-10-len(repo)))

		outputCompletes := checkGitLog([]string{"completes", "completed", "complete"}, story, dir, storyResetTime)
		outputFixes := checkGitLog([]string{"fixes", "fixed", "fix"}, story, dir, storyResetTime)
		outputFinishes := checkGitLog([]string{"finishes", "finished", "finish"}, story, dir, storyResetTime)

		if len(outputCompletes) > 0 || len(outputFixes) > 0 || len(outputFinishes) > 0 {
			sayf(colorstring.Color("[green]DELIVERING\n"))
			if comment == "" {
				client.DeliverStory(story.ID)
			} else {
				client.DeliverStoryWithComment(story.ID, comment)
			}
		} else {
			sayf(colorstring.Color("  [yellow]SKIPPING\n"))
		}
	}

	sayf("\n")
}

var resettingStoryStates map[string]bool = map[string]bool{
	"delivered": true,
	"rejected":  true,
	"accepted":  true,
}

func getStoryResetTime(client tracker.ProjectClient, story tracker.Story) (time.Time, error) {
	activityList, err := client.StoryActivity(story.ID, tracker.ActivityQuery{})
	if err != nil {
		return time.Time{}, err
	}
	for _, activity := range activityList {
		if activity.Kind == "story_update_activity" && resettingStoryStates[activity.Highlight] {
			return activity.OccurredAt, nil
		}
	}
	return time.Time{}, nil
}

func checkGitLog(verbs []string, story tracker.Story, dir string, after time.Time) []byte {
	afterArg := after.Format(time.RFC3339)
	verbsRegexp := fmt.Sprintf("(%s)", strings.Join(verbs, "|"))
	command := exec.Command("git", "log", "-i", "--extended-regexp", "--grep", fmt.Sprintf("%s\\s+(#[0-9]+,\\s+)*#%d", verbsRegexp, story.ID), "--after", afterArg)
	command.Dir = dir

	output, err := command.CombinedOutput()
	if err != nil {
		sayf("git logging failed for story: %d: %s\n", story.ID, err)

		return nil
	}

	command = exec.Command("git", "log", "-i", "--extended-regexp", "--grep", fmt.Sprintf("#%d(,\\s+#[0-9]+)*\\s+%s", story.ID, verbsRegexp), "--after", afterArg)
	command.Dir = dir

	output2, err := command.CombinedOutput()
	if err != nil {
		sayf("git logging failed for story: %d: %s\n", story.ID, err)

		return nil
	}

	return append(output, output2...)
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
