package smartling

import "fmt"

const (
	endpointFileStatusExtended = "/files-api/v2/projects/%s/locales/%s/file/status"
)

// GetFileStatus returns extended file status for a file for a particular locale
func (client *Client) GetFileStatusExtended(
	projectID string,
	fileURI string,
	localeID string,
) (*FileStatusExtended, error) {
	var status FileStatusExtended

	_, _, err := client.GetJSON(
		fmt.Sprintf(endpointFileStatusExtended, projectID, localeID),
		FileURIRequest{fileURI}.GetQuery(),
		&status,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get files list: %s", err,
		)
	}

	return &status, nil
}
