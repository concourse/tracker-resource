# Tracker Resource

An output only resource that will deliver finished [Pivotal Tracker][tracker] stories that are linked in recent git commits.

[tracker]: https://www.pivotaltracker.com

## Usage

When you use the [Pivotal Tracker syntax for finishing a story or fixing a bug](https://www.pivotaltracker.com/help/articles/githubs_service_hook_for_tracker/#formatting-your-commits) this resource will detect it in your commit history and deliver the appropriate stories.

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

You'll need a separate resource for each Tracker project.

### Build

``` yaml
- name: deploy
  plan:
  - get: git-repo-path
  - ...
  - put: tracker-output
    params:
      repos:
      - git-repo-path
```

#### Out Parameters

* `repos`: *Required.* Paths to the git repositories which will contain the delivering commits.

* `comment`: *Optional.* A file containing a comment to leave on any delivered stories.

## Development

### Running the tests

The tests have been embedded with the `Dockerfile`; ensuring that the testing
environment is consistent across any `docker` enabled platform. When the docker
image builds, the test are run inside the docker container, on failure they
will stop the build.

Run the tests with the following commands for both `alpine` and `ubuntu` images:

```sh
docker build -t tracker-resource -f dockerfiles/alpine/Dockerfile .
docker build -t tracker-resource -f dockerfiles/ubuntu/Dockerfile .
```
