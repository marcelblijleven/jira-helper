package pkg

import (
	"github.com/stretchr/testify/assert"
	"log"
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
	assert.Equal(t, []string{"MB-1337"}, result)
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

	err = AssignVersions(releaseBody, "My first version", jiraClient)
	assert.NoError(t, err)
	assert.Equal(t, 1, mockClient.CalledTimes)
	assert.Equal(t, "{\"update\":{\"fixVersions\":[{\"add\":{\"name\":\"My first version\"}}]}}", mockClient.CalledWith[0])
}

func TestAssignVersions_multipleTickets(t *testing.T) {
	mockClient := NewMockHttpClient(t, 201)
	jiraClient, err := NewJiraClient("https://test.nu", "marcel@test.nu", "c0ffee", mockClient)

	if err != nil {
		log.Fatalln(err)
	}

	releaseBody := `This is an automated release.
For changes in this version, see the changelog

Merge commit that triggered this release: feat: marcel introduces c0ffee (MB-1337, MB-1338)`

	err = AssignVersions(releaseBody, "My first version", jiraClient)
	assert.NoError(t, err)
	assert.Equal(t, 2, mockClient.CalledTimes)
	assert.Equal(t, []string{
		"{\"update\":{\"fixVersions\":[{\"add\":{\"name\":\"My first version\"}}]}}",
		"{\"update\":{\"fixVersions\":[{\"add\":{\"name\":\"My first version\"}}]}}",
	}, mockClient.CalledWith)
}
