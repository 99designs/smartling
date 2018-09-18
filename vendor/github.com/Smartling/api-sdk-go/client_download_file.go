package smartling

import (
	"fmt"
	"io"
)

const (
	endpointDownloadFile = "/files-api/v2/projects/%s/file"
)

// DownloadFile downloads original file from project.
func (client *Client) DownloadFile(
	projectID string,
	uri string,
) (io.Reader, error) {
	reader, _, err := client.Get(
		fmt.Sprintf(endpointDownloadFile, projectID),
		FileURIRequest{FileURI: uri}.GetQuery(),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to download original file: %s", err,
		)
	}

	return reader, nil
}
