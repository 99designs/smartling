package smartling

import (
	"fmt"
	"io"
)

const (
	endpointDownloadTranslation = "/files-api/v2/projects/%s/locales/%s/file"
)

// DownloadTranslation downloads specified translated file for specified
// locale.  Check FileDownloadRequest for more options.
func (client *Client) DownloadTranslation(
	projectID string,
	localeID string,
	request FileDownloadRequest,
) (io.Reader, error) {
	reader, _, err := client.Get(
		fmt.Sprintf(endpointDownloadTranslation, projectID, localeID),
		request.GetQuery(),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to download translated file: %s", err,
		)
	}

	return reader, nil
}
