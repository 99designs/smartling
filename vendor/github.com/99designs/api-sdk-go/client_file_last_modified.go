package smartling

import "fmt"

const (
	endpointFilesLastModified = "/files-api/v2/projects/%s/file/last-modified"
)

func (client *Client) LastModified(
	projectID string,
	request FileLastModifiedRequest,
) (*FileLastModifiedLocales, error) {
	var lastModifiedLocales FileLastModifiedLocales

	_, _, err := client.GetJSON(
		fmt.Sprintf(endpointFilesLastModified, projectID),
		request.GetQuery(),
		&lastModifiedLocales,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get last modified: %s", err,
		)
	}

	return &lastModifiedLocales, nil
}
