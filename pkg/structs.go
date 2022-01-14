package pkg

// releaseRequestBody represents the Jira create release API request body
type releaseRequestBody struct {
	Name        string `json:"name"`
	Released    bool   `json:"released"`
	ReleaseDate string `json:"releaseDate"`
	Project     string `json:"project"`
}

// assignRequestBody represents the Jira assign fixVersion API request body
type assignRequestBody struct {
	Update update `json:"update"`
}

type update struct {
	FixVersions []fixVersion `json:"fixVersions"`
}

type fixVersion struct {
	Add addFixVersion `json:"add"`
}

type addFixVersion struct {
	Name string `json:"name"`
}

type JiraError struct {
	ErrorMessages []interface{} `json:"errorMessages"`
	Errors        struct {
		Name string `json:"name"`
	} `json:"errors"`
}

// createResponse represents the response from the Jira api when calling the create fix version endpoint
type createResponse struct {
	Self            string `json:"self"`
	Id              string `json:"id"`
	Name            string `json:"name"`
	Archived        bool   `json:"archived"`
	Released        bool   `json:"released"`
	ReleaseDate     string `json:"releaseDate"`
	UserReleaseDate string `json:"userReleaseDate"`
	ProjectId       int    `json:"projectId"`
}
