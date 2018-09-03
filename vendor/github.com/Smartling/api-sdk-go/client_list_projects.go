package smartling

import (
	"fmt"
)

const (
	endpointProjectsList   = "/accounts-api/v2/accounts/%v/projects"
	endpointProjectDetails = "/projects-api/v2/projects/%v"
)

// ProjectsList represents projects list under specified account.
type ProjectsList struct {
	// TotalCount represents total count of projects.
	TotalCount int64

	// Items contains projects list by specified request.
	Items []Project
}

// Project represents detailed project information.
type Project struct {
	// ProjectID is a unique project ID.
	ProjectID string

	// ProjectName is a human-friendly project name.
	ProjectName string

	// AccountUID is undocumented by Smartling API.
	AccountUID string

	// SourceLocaleID represents source locale ID for project.
	SourceLocaleID string

	// SourceLocaleDescription describes project's locale.
	SourceLocaleDescription string

	// Archived will be true if project is archived.
	Archived bool
}

// ProjectDetails extends Project type to contain target locales list.
type ProjectDetails struct {
	Project

	// TargetLocales represents target locales list.
	TargetLocales []Locale
}

// Locale represents locale for translation.
type Locale struct {
	// LocaleID is a unique locale ID.
	LocaleID string

	// Description describes locale.
	Description string

	// Enabled is a flag that represents is locale enabled or not.
	Enabled bool
}

// ListProjects returns projects in specified account matching specified
// request.
func (client *Client) ListProjects(
	accountID string,
	request ProjectsListRequest,
) (*ProjectsList, error) {
	var list ProjectsList

	_, _, err := client.GetJSON(
		fmt.Sprintf(endpointProjectsList, accountID),
		request.GetQuery(),
		&list,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get projects list: %s", err,
		)
	}

	return &list, nil
}

// GetProjectDetails returns project details for specified project.
func (client *Client) GetProjectDetails(
	projectID string,
) (*ProjectDetails, error) {
	var details ProjectDetails

	_, _, err := client.GetJSON(
		fmt.Sprintf(endpointProjectDetails, projectID),
		nil,
		&details,
	)
	if err != nil {
		if _, ok := err.(NotFoundError); ok {
			return nil, err
		}

		return nil, fmt.Errorf(
			"failed to get project details: %s", err,
		)
	}

	return &details, nil
}
