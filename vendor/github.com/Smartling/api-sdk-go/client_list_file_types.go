package smartling

import "fmt"

// ListFileTypes returns returns file types list from specified project.
func (client *Client) ListFileTypes(
	projectID string,
) ([]FileType, error) {
	var result struct {
		Items []FileType
	}

	_, _, err := client.Get(
		fmt.Sprintf(endpointFileTypes, projectID),
		nil,
		&result,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get file types: %s", err,
		)
	}

	return result.Items, nil
}
