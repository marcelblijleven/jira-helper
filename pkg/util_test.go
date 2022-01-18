package pkg

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func Test_extractIssuesFromText(t *testing.T) {
	body := `This is an automated release.
For changes in this version, see the changelog

Merge commit that triggered this release: feat: marcel introduces c0ffee (MB-1337, HB-1338)`

	result := extractIssuesFromText(body)
	assert.Equal(t, []string{"MB-1337", "HB-1338"}, result)
}

func Test_extractIssuesFromText_duplicateTickets(t *testing.T) {
	body := `This is an automated release.
For changes in this version, see the changelog

Merge commit that triggered this release: feat: marcel introduces c0ffee (MB-1337, MB-1337)`

	result := extractIssuesFromText(body)
	assert.Equal(t, []string{"MB-1337", "MB-1337"}, result)
}

func Test_removeDuplicates(t *testing.T) {
	type args struct {
		items []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "no duplicates",
			args: args{items: []string{"I", "like", "c0ffee"}},
			want: []string{"I", "like", "c0ffee"},
		},
		{
			name: "nil slice",
			args: args{items: nil},
			want: []string(nil),
		},
		{
			name: "empty slice",
			args: args{items: []string{}},
			want: []string{},
		},
		{
			name: "slice with duplicates",
			args: args{items: []string{"a", "b", "a"}},
			want: []string{"a", "b"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := removeDuplicates(tt.args.items)
			assert.Equal(t, tt.want, got, "removeDuplicates() = %v, want %v", got, tt.want)
		})
	}
}

func TestAssignVersions(t *testing.T) {
	mockClient := NewMockHttpClient(t, 201)
	jiraClient, err := NewJiraClient("https://test.nu", "marcel@test.nu", "c0ffee", mockClient)

	if err != nil {
		log.Fatalln(err)
	}

	releaseBody := `This is an automated release.
For changes in this version, see the changelog

Merge commit that triggered this release: feat: marcel introduces c0ffee (MB-1337)`

	err = AssignVersions(releaseBody, "My first version", jiraClient, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, mockClient.CalledTimes)
	assert.Equal(t, "{\"update\":{\"fixVersions\":[{\"add\":{\"name\":\"My first version\"}}]}}", mockClient.CalledWith[0])
}

func TestAssignVersions_multipleIssues(t *testing.T) {
	mockClient := NewMockHttpClient(t, 201)
	jiraClient, err := NewJiraClient("https://test.nu", "marcel@test.nu", "c0ffee", mockClient)

	if err != nil {
		log.Fatalln(err)
	}

	releaseBody := `This is an automated release.
For changes in this version, see the changelog

Merge commit that triggered this release: feat: marcel introduces c0ffee (MB-1337, MB-1338)`

	err = AssignVersions(releaseBody, "My first version", jiraClient, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, 2, mockClient.CalledTimes)
	assert.Equal(t, []string{
		"{\"update\":{\"fixVersions\":[{\"add\":{\"name\":\"My first version\"}}]}}",
		"{\"update\":{\"fixVersions\":[{\"add\":{\"name\":\"My first version\"}}]}}",
	}, mockClient.CalledWith)
}

func TestAssignVersions_singleIssue(t *testing.T) {
	mockClient := NewMockHttpClient(t, 201)
	jiraClient, err := NewJiraClient("https://test.nu", "marcel@test.nu", "c0ffee", mockClient)

	if err != nil {
		log.Fatalln(err)
	}

	err = AssignVersions("", "My first version", jiraClient, []string{"MB-1234"}, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, mockClient.CalledTimes)
	assert.Equal(t, []string{
		"{\"update\":{\"fixVersions\":[{\"add\":{\"name\":\"My first version\"}}]}}",
	}, mockClient.CalledWith)
}

func TestAssignVersions_duplicateIssue(t *testing.T) {
	mockClient := NewMockHttpClient(t, 201)
	jiraClient, err := NewJiraClient("https://test.nu", "marcel@test.nu", "c0ffee", mockClient)

	if err != nil {
		log.Fatalln(err)
	}

	err = AssignVersions("", "My first version", jiraClient, []string{"MB-1234", "MB-1234"}, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, mockClient.CalledTimes)
	assert.Equal(t, []string{
		"{\"update\":{\"fixVersions\":[{\"add\":{\"name\":\"My first version\"}}]}}",
	}, mockClient.CalledWith)
}

func TestAssignVersions_multipleIssuesFromBodyAndSeparate(t *testing.T) {
	mockClient := NewMockHttpClient(t, 201)
	jiraClient, err := NewJiraClient("https://test.nu", "marcel@test.nu", "c0ffee", mockClient)

	if err != nil {
		log.Fatalln(err)
	}

	releaseBody := `This is an automated release.
For changes in this version, see the changelog

Merge commit that triggered this release: feat: marcel introduces c0ffee (MB-1337, MB-1338)`

	err = AssignVersions(releaseBody, "My first version", jiraClient, []string{"MB-1339", "MB-1340"}, nil)
	assert.NoError(t, err)
	assert.Equal(t, 4, mockClient.CalledTimes)
	assert.Equal(t, []string{
		"{\"update\":{\"fixVersions\":[{\"add\":{\"name\":\"My first version\"}}]}}",
		"{\"update\":{\"fixVersions\":[{\"add\":{\"name\":\"My first version\"}}]}}",
		"{\"update\":{\"fixVersions\":[{\"add\":{\"name\":\"My first version\"}}]}}",
		"{\"update\":{\"fixVersions\":[{\"add\":{\"name\":\"My first version\"}}]}}",
	}, mockClient.CalledWith)
}

func TestAssignVersions_multipleIssuesFromBodyAndSeparate_withDuplicates(t *testing.T) {
	mockClient := NewMockHttpClient(t, 201)
	jiraClient, err := NewJiraClient("https://test.nu", "marcel@test.nu", "c0ffee", mockClient)

	if err != nil {
		log.Fatalln(err)
	}

	releaseBody := `This is an automated release.
For changes in this version, see the changelog

Merge commit that triggered this release: feat: marcel introduces c0ffee (MB-1337, MB-1338)`

	err = AssignVersions(releaseBody, "My first version", jiraClient, []string{"MB-1337", "MB-1338"}, nil)
	assert.NoError(t, err)
	assert.Equal(t, 2, mockClient.CalledTimes)
	assert.Equal(t, []string{
		"{\"update\":{\"fixVersions\":[{\"add\":{\"name\":\"My first version\"}}]}}",
		"{\"update\":{\"fixVersions\":[{\"add\":{\"name\":\"My first version\"}}]}}",
	}, mockClient.CalledWith)
}

func Test_handleJiraError(t *testing.T) {
	body := bytes.NewReader([]byte("{\"errorMessages\":[],\"errors\":{\"name\":\"A version with this name already exists in this project.\"}}"))
	res := &http.Response{Status: "Bad request", StatusCode: http.StatusBadRequest, Body: ioutil.NopCloser(body)}
	err := handleJiraError(res)
	assert.EqualError(t, err, "request unsuccessful (Bad request): A version with this name already exists in this project.")
}

func Test_handleJiraError_unreadableResponse(t *testing.T) {
	body := bytes.NewReader([]byte("test body"))
	res := &http.Response{Status: "Bad request", StatusCode: http.StatusBadRequest, Body: ioutil.NopCloser(body)}
	err := handleJiraError(res)
	assert.EqualError(t, err, "request unsuccessful (Bad request), could not read response: invalid character 'e' in literal true (expecting 'r')")
}

func Test_filterSlice(t *testing.T) {
	type args struct {
		main   []string
		filter []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "empty filter",
			args: args{
				main:   []string{"a", "b", "c"},
				filter: []string{},
			},
			want: []string{"a", "b", "c"},
		},
		{
			name: "nil filter",
			args: args{
				main:   []string{"a", "b", "c"},
				filter: nil,
			},
			want: []string{"a", "b", "c"},
		},
		{
			name: "filter some",
			args: args{
				main:   []string{"a", "b", "c"},
				filter: []string{"b"},
			},
			want: []string{"a", "c"},
		},
		{
			name: "filter all",
			args: args{
				main:   []string{"a", "b", "c"},
				filter: []string{"a", "b", "c"},
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := filterSlice(tt.args.main, tt.args.filter)
			assert.Equal(t, tt.want, got, "filterSlice() = %v, want %v", got, tt.want)
		})
	}
}
