package pkg

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

type MockHttpClient struct {
	t             *testing.T
	CalledMethod  string
	CalledWith    []string
	CalledTimes   int
	CalledHeaders http.Header
	statusCode    int
}

func NewMockHttpClient(t *testing.T, statusCode int) *MockHttpClient {
	return &MockHttpClient{t: t, statusCode: statusCode, CalledTimes: 0, CalledWith: []string{}}
}

func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	m.CalledTimes += 1
	data, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()

	if err != nil {
		m.t.Fatal("error occurred while doing mock request")
	}

	m.CalledWith = append(m.CalledWith, string(data))
	m.CalledMethod = req.Method
	m.CalledHeaders = req.Header

	if m.statusCode == http.StatusBadRequest {
		return &http.Response{StatusCode: http.StatusBadRequest, Status: "bad request"}, nil
	}

	return &http.Response{StatusCode: m.statusCode}, nil
}

func TestNewJiraClient(t *testing.T) {
	m := NewMockHttpClient(t, 200)
	parsedHost, _ := url.Parse("https://test.nu")

	c, err := NewJiraClient("https://test.nu", "marcel@test.nl", "c0ffee", m)

	assert.NoError(t, err)

	expected := &JiraClient{
		host:           parsedHost,
		httpClient:     m,
		authentication: nil,
	}

	expected.authentication = &authenticationService{
		client: expected,
		email:  "marcel@test.nl",
		token:  "c0ffee",
	}

	assert.Equal(t, expected, c)
}

func TestNewJiraClient_missingHost(t *testing.T) {
	c, err := NewJiraClient("", "marcel@test.nl", "c0ffee", nil)
	assert.Nil(t, c)
	assert.EqualError(t, err, "could not create jira client: hostname cannot be empty")
}

func TestNewJiraClient_missingEmail(t *testing.T) {
	c, err := NewJiraClient("https://test.nu", "", "c0ffee", nil)
	assert.Nil(t, c)
	assert.EqualError(t, err, "could not create jira client: email cannot be empty")
}

func TestNewJiraClient_missingToken(t *testing.T) {
	c, err := NewJiraClient("https://test.nu", "marcel@test.nl", "", nil)
	assert.Nil(t, c)
	assert.EqualError(t, err, "could not create jira client: token cannot be empty")
}

func TestJiraClient_CreateFixVersion(t *testing.T) {
	mockClient := NewMockHttpClient(t, 201)
	jiraClient, err := NewJiraClient("https://test.nu", "marcel@test.nl", "c0ffee", mockClient)

	if err != nil {
		t.Fatal(err)
	}

	err = jiraClient.CreateFixVersion("test version", "MB")
	assert.NoError(t, err)
	assert.Equal(t, 1, mockClient.CalledTimes)
	assert.Equal(t, fmt.Sprintf("{\"name\":\"test version\",\"released\":true,\"releaseDate\":\"%v\",\"project\":\"MB\"}", getDateString()), mockClient.CalledWith[0])
	assert.Equal(t, http.MethodPost, mockClient.CalledMethod)
	assert.Equal(t, "Basic bWFyY2VsQHRlc3Qubmw6YzBmZmVl", mockClient.CalledHeaders.Get("Authorization"))
	assert.Equal(t, "application/json", mockClient.CalledHeaders.Get("Content-Type"))
}

func TestJiraClient_CreateFixVersion_non20X(t *testing.T) {
	mockClient := NewMockHttpClient(t, 400)
	jiraClient, err := NewJiraClient("https://test.nu", "marcel@test.nl", "c0ffee", mockClient)

	if err != nil {
		t.Fatal(err)
	}

	err = jiraClient.CreateFixVersion("test version", "MB")
	assert.Equal(t, 1, mockClient.CalledTimes)
	assert.EqualError(t, err, "create fix version unsuccessful: bad request")
}

func TestJiraClient_AssignVersion(t *testing.T) {
	mockClient := NewMockHttpClient(t, 201)
	jiraClient, err := NewJiraClient("https://test.nu", "marcel@test.nl", "c0ffee", mockClient)

	if err != nil {
		t.Fatal(err)
	}

	err = jiraClient.AssignVersion("MB-1337", "My first release")
	assert.NoError(t, err)
	assert.Equal(t, 1, mockClient.CalledTimes)
	assert.Equal(t, "{\"update\":{\"fixVersions\":[{\"add\":{\"name\":\"My first release\"}}]}}", mockClient.CalledWith[0])
	assert.Equal(t, http.MethodPut, mockClient.CalledMethod)
	assert.Equal(t, "Basic bWFyY2VsQHRlc3Qubmw6YzBmZmVl", mockClient.CalledHeaders.Get("Authorization"))
	assert.Equal(t, "application/json", mockClient.CalledHeaders.Get("Content-Type"))
}

func TestJiraClient_AssignVersion_non20X(t *testing.T) {
	mockClient := NewMockHttpClient(t, 400)
	jiraClient, err := NewJiraClient("https://test.nu", "marcel@test.nl", "c0ffee", mockClient)

	if err != nil {
		t.Fatal(err)
	}

	err = jiraClient.AssignVersion("MB-1337", "My first release")
	assert.EqualError(t, err, "assign version unsuccessful: bad request")
}
