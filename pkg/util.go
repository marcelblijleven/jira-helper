package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

// getDateString returns the current date in YYYY-MM-DD format
func getDateString() string {
	return time.Now().Format("2006-01-02")
}

// newReleaseRequestBody creates a release request body with the provided version name and project id
func newReleaseRequestBody(versionName, projectID string) (*releaseRequestBody, error) {
	if versionName == "" {
		return nil, errors.New("version versionName cannot be empty")
	}

	if projectID == "" {
		return nil, errors.New("project ID cannot be empty")
	}

	return &releaseRequestBody{
		Name:        versionName,
		Released:    true,
		ReleaseDate: getDateString(),
		Project:     projectID,
	}, nil
}

// newAssignRequestBody creates an assign fixVersion request body with the provided version
func newAssignRequestBody(version string) (*assignRequestBody, error) {
	if version == "" {
		return nil, errors.New("version cannot be empty")
	}

	f := fixVersion{Add: addFixVersion{Name: version}}
	b := &assignRequestBody{Update: update{FixVersions: []fixVersion{f}}}
	return b, nil
}

// removeDuplicates can be used to filter out duplicates from the provided slice of string
func removeDuplicates(items []string) []string {
	if items == nil || len(items) == 0 {
		return items
	}

	var filtered []string
	check := make(map[string]bool)

	for _, item := range items {
		if _, ok := check[item]; !ok {
			check[item] = true
			filtered = append(filtered, item)
		}
	}

	return filtered
}

// extractIssuesFromText gathers all issue numbers from the provided text
func extractIssuesFromText(text string) []string {
	r := regexp.MustCompile("[A-Z]+-[0-9]+")
	return r.FindAllString(text, -1)
}

// AssignVersions extracts the issues from  the provided release body and calls the AssignVersion endpoint of the
// jira client.
func AssignVersions(releaseBody, version string, client *JiraClient, issues ...string) error {
	issues = append(issues, extractIssuesFromText(releaseBody)...)
	issues = removeDuplicates(issues)

	for _, issue := range issues {
		if err := client.AssignVersion(issue, version); err != nil {
			return fmt.Errorf("error occurred while assign version to issue %s: %w", issue, err)
		}
	}

	return nil
}

// handleJiraError retrieves and formats the error from the Jira api response
func handleJiraError(res *http.Response) error {
	var jiraError JiraError
	data, readErr := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if readErr != nil {
		return fmt.Errorf("request unsuccessful (%s), could not read response: %w", res.Status, readErr)
	}

	if unmarshallErr := json.Unmarshal(data, &jiraError); unmarshallErr != nil {
		return fmt.Errorf("request unsuccessful (%s), could not read response: %w", res.Status, unmarshallErr)
	}

	return fmt.Errorf("request unsuccessful (%s): %s", res.Status, jiraError.Errors.Name)
}
