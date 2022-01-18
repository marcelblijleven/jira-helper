package cmd

var (
	user    string
	host    string
	token   string
	project string
	version string
	body    string
	issues  []string
	filter  []string
)

const (
	userFlagName  = "user"
	userShorthand = "u"
	userUsage     = "User (email) for authenticating against the Jira API"

	hostFlagName  = "host"
	hostShorthand = "s"
	hostUsage     = "Host of the Jira API. If the host URL contains a scheme (e.g. https), you must include it"

	tokenFlagName  = "token"
	tokenShorthand = "t"
	tokenUsage     = "Token used to authenticate against the Jira API"

	projectFlagName  = "project"
	projectShorthand = "p"
	projectUsage     = "Project key of the Jira project, e.g. MB"

	versionFlagName  = "version"
	versionShorthand = "v"
	versionUsage     = "Name of the version"

	bodyFlagName  = "releaseBody"
	bodyShorthand = "b"
	bodyUsage     = "The body of text which contains Jira issues, e.g. a GitHub release body"

	issuesFlagName  = "issues"
	issuesShorthand = "i"
	issuesUsage     = "The issues you want to assign to release to, can be a single issue or comma separated"

	filterFlagName  = "filter"
	filterShorthand = "f"
	filterUsage     = "The filter flag allows you to ignore issues when assigning a release"
)
