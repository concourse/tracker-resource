# Tracker Resource

An output only resource that will deliver finished [Pivotal Tracker][tracker] stories that are linked in recent git commits.

[tracker]: https://www.pivotaltracker.com

## Configuration

### Resource

``` yaml
- name: tracker-output
  type: tracker
  source:
    token: TRACKER_API_TOKEN
    project_id: "TRACKER_PROJECT_ID"
    tracker_url: https://www.pivotaltracker.com
```

Replace `TRACKER_API_TOKEN` and `TRACKER_PROJECT_ID` with your API token (which can be found on your profile page) and your project ID (which can be found in the URL of your project). Make sure that your project ID is a string because it will converted to JSON when given to the resource and JSON doesn't like integers.

You'll need a seperate resource for each Tracker project.

### Build

``` yaml
- name: deploy
  plan:
  - get: git-repo-path
  - ...
  - put: tracker
    params:
      repos:
      - git-repo-path
```

#### Out Parameters

* `repos`: *Required.* Paths to the git repositories which will contain the delivering commits.

* `comment`: *Optional.* A file containing a comment to leave on any delivered stories. If not specified, this comment defaults to "Delivered by Concourse."
