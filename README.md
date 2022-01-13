# Jira Helper
Helper tool to interact with Jira from CI/CD scripts. Its main purpose is to create and assign version
based on GitHub releases to Jira tickets.

## GitHub actions example
The following can be used to trigger a release (fixVersion) in Jira whenever a GitHub release is created. It will use the body of the
release and search for any issue numbers in it and automatically assign the newly created release to them.
```yaml
name: Release

on:
  release:
    types:
      - released

jobs:
  release-to-jira:
    runs-on: ubuntu-latest
    env:
      IMAGE_NAME: ghcr.io/marcelblijleven/jira-helper:latest
    steps:
      - name: create release in Jira
        run: docker run -i --rm ${{ env.IMAGE_NAME }} createRelease -u marcel@test.nu -s https://your-jira.address.nl -p MB -t=${{ secrets.API_TOKEN }} -v "${{ github.event.release.name }}"
      - name: assign release to Jira tickets
        run: docker run -i --rm ${{ env.IMAGE_NAME }} assignVersion -u marcel@test.nu -s https://your-jira.address.nl -p MB -t=${{ secrets.API_TOKEN }} -v "${{ github.event.release.name }}" -b "${{ github.event.release.body }}"

```

## CLI Usage
```
Usage:
  jira-helper [command]

Available Commands:
  assignVersion Assigns a version to all provided issues in the release body
  completion    Generate the autocompletion script for the specified shell
  createRelease Create a fix version in Jira
  help          Help about any command

Flags:
  -h, --help   help for jira-helper

```

### Assign version
Assigns a version to all provided issues. The issue numbers are retrieved from
the provided release body.

```
Usage:
  jira-helper assignVersion [flags]

Flags:
  -h, --help                 help for assignVersion
  -s, --host string          host of the Jira API
  -p, --project string       Abbreviation of the Jira project, e.g. GGWM
  -b, --releaseBody string   The body of the Github release
  -t, --token string         Token used to authenticate against the Jira API
  -u, --user string          user used for authenticating against the Jira API
  -v, --version string       Version name
```

### Create release
Create a fix version in Jira for the project with the provided name.

The release state of the fix version will be set to "released" and the day will be set to
today.

```
Usage:
  jira-helper createRelease [flags]

Flags:
  -h, --help             help for createRelease
  -s, --host string      host of the Jira API
  -p, --project string   Abbreviation of the Jira project, e.g. GGWM
  -t, --token string     Token used to authenticate against the Jira API
  -u, --user string      user used for authenticating against the Jira API
  -v, --version string   Version name

```
