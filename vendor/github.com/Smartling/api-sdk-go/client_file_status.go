package smartling

import "fmt"

const (
	endpointFileStatus = "/files-api/v2/projects/%s/file/status"
)

// GetFileStatus returns file status.
func (client *Client) GetFileStatus(
	projectID string,
	fileURI string,
) (*FileStatus, error) {
	var status FileStatus

	_, _, err := client.GetJSON(
		fmt.Sprintf(endpointFileStatus, projectID),
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
