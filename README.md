# Jira Helper
Helper tool to interact with Jira from CI/CD scripts. Its main purpose is to create and assign version
based on GitHub releases to Jira tickets.

You can provide comma separated issue numbers or a body of text which contains issue numbers and the jira-helper
will automatically assign the version to those issues in Jira.

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
  assignRelease   Assigns a version to all provided issues in the release body
  completion      Generate the autocompletion script for the specified shell
  createAndAssign Creates a fix version in Jira and assigns it to the issues
  createRelease   Create a fix version in Jira
  help            Help about any command

Flags:
  -h, --help             help for jira-helper
  -s, --host string      Host of the Jira API. If the host URL contains a scheme (e.g. https), you must include it
  -p, --project string   Project key of the Jira project, e.g. MB
  -t, --token string     Token used to authenticate against the Jira API
  -u, --user string      User (email) for authenticating against the Jira API
  -v, --version string   Name of the version

```

### Assign release
Assigns a version to all provided issues. The issue numbers are retrieved from
the provided release body.

```
Usage:
jira-helper assignRelease [flags]

Aliases:
assignRelease, assignVersion

Flags:
-f, --filter strings       The filter flag allows you to ignore issues when assigning a release
-h, --help                 help for assignRelease
-i, --issues strings       The issues you want to assign to release to, can be a single issue or comma separated
-b, --releaseBody string   The body of text which contains Jira issues, e.g. a GitHub release body

Global Flags:
-s, --host string      Host of the Jira API. If the host URL contains a scheme (e.g. https), you must include it
-p, --project string   Project key of the Jira project, e.g. MB
-t, --token string     Token used to authenticate against the Jira API
-u, --user string      User (email) for authenticating against the Jira API
-v, --version string   Name of the version
```

### Create release
Create a fix version in Jira for the project with the provided name.

The release state of the fix version will be set to "released" and the day will be set to
today.

```
Usage:
jira-helper createRelease [flags]

Aliases:
createRelease, createVersion

Flags:
-h, --help   help for createRelease

Global Flags:
-s, --host string      Host of the Jira API. If the host URL contains a scheme (e.g. https), you must include it
-p, --project string   Project key of the Jira project, e.g. MB
-t, --token string     Token used to authenticate against the Jira API
-u, --user string      User (email) for authenticating against the Jira API
-v, --version string   Name of the version

```

### Create and assign
Creates a fix version in Jira and assigns it to the provided issues.

The release state of the fix version will be set to "released" and the day will be set to
today.

```
Usage:
jira-helper createAndAssign [flags]

Flags:
-f, --filter strings       The filter flag allows you to ignore issues when assigning a release
-h, --help                 help for createAndAssign
-i, --issues strings       The issues you want to assign to release to, can be a single issue or comma separated
-b, --releaseBody string   The body of text which contains Jira issues, e.g. a GitHub release body

Global Flags:
-s, --host string      Host of the Jira API. If the host URL contains a scheme (e.g. https), you must include it
-p, --project string   Project key of the Jira project, e.g. MB
-t, --token string     Token used to authenticate against the Jira API
-u, --user string      User (email) for authenticating against the Jira API
-v, --version string   Name of the version

```
