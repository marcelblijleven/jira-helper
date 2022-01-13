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
