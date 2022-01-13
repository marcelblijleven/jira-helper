package pkg

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const (
	apiEndpoint = "/rest/api/latest"
)

// JiraClient allows the user to interact with the Jira API to create and assign fixVersions
type JiraClient struct {
	host           *url.URL
	httpClient     HttpClient
	authentication *authenticationService
}

// HttpClient is the http client interface used by the Jira client
type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// NewJiraClient creates a new JiraClient with the provided values
func NewJiraClient(host, email, token string, httpClient HttpClient) (*JiraClient, error) {
	if host == "" {
		return nil, errors.New("could not create jira client: hostname cannot be empty")
	}

	if email == "" {
		return nil, errors.New("could not create jira client: email cannot be empty")
	}

	if token == "" {
		return nil, errors.New("could not create jira client: token cannot be empty")
	}

	u, err := url.Parse(host)

	if err != nil {
		return nil, fmt.Errorf("could not create jira client, invalid host provided: %w", err)
	}

	client := &JiraClient{host: u, httpClient: httpClient}
	client.authentication = &authenticationService{
		client: client,
		email:  email,
		token:  token,
	}
	return client, nil
}

// createRequest creates a request with the provided method and body
func (c *JiraClient) createRequest(method, endpoint string, body interface{}) (*http.Request, error) {
	e, err := url.Parse(endpoint)

	if err != nil {
		return nil, fmt.Errorf("could not parse endpoint: %w", err)
	}

	data, err := json.Marshal(body)

	if err != nil {
		return nil, fmt.Errorf("could not marshall provided body to json: %w", err)
	}

	u := c.host.ResolveReference(e).String()

	req, err := http.NewRequest(method, u, bytes.NewReader(data))

	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	c.authentication.setBasicAuth(req)
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

// AssignVersion calls the issue endpoint to add a fixVersion to the issue
func (c *JiraClient) AssignVersion(issue, version string) error {
	endpoint := fmt.Sprintf("%s/issue/%s", apiEndpoint, issue)
	body, err := newAssignRequestBody(version)

	if err != nil {
		return fmt.Errorf("could not create assign version request body: %w", err)
	}

	req, err := c.createRequest(http.MethodPut, endpoint, body)

	if err != nil {
		return err
	}

	res, err := c.httpClient.Do(req)

	if err != nil {
		return fmt.Errorf("could not do request: %w", err)
	}

	if res.StatusCode/100 != 2 {
		return fmt.Errorf("assign version unsuccessful: %v", res.Status)
	}

	return nil
}

// CreateFixVersion calls the version endpoint to add a fixVersion to the provided project
func (c *JiraClient) CreateFixVersion(name, project string) error {
	endpoint := apiEndpoint + "/version"
	body, err := newReleaseRequestBody(name, project)

	if err != nil {
		return fmt.Errorf("could not create new release request body: %w", err)
	}

	req, err := c.createRequest(http.MethodPost, endpoint, body)

	if err != nil {
		return err
	}

	res, err := c.httpClient.Do(req)

	if err != nil {
		return fmt.Errorf("could not do request: %w", err)
	}

	if res.StatusCode/100 != 2 {
		return fmt.Errorf("create fix version unsuccessful: %v", res.Status)
	}

	return nil
}
