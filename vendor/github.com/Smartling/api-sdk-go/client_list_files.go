package smartling

import (
	"fmt"
)

const (
	endpointFilesList = "/files-api/v2/projects/%s/files/list"
	endpointFileTypes = "/files-api/v2/projects/%s/file-types"
)

// FilesList represents file list reply from Smartling APIa.
type FilesList struct {
	// TotalCount is a total files count.
	TotalCount int

	// Items contains all files matched by request.
	Items []File
}

// ListFiles returns files list from specified project by specified request.
// Returned result is paginated, so check out TotalCount struct field in the
// reply. API can return only 500 files at once.
func (client *Client) ListFiles(
	projectID string,
	request FilesListRequest,
) (*FilesList, error) {
	var list FilesList

	_, _, err := client.GetJSON(
		fmt.Sprintf(endpointFilesList, projectID),
		request.GetQuery(),
		&list,
	)
	if err != nil {
		if _, ok := err.(NotFoundError); ok {
			return nil, err
		}

		return nil, fmt.Errorf(
			"failed to get files list: %s", err,
		)
	}

	return &list, nil
}

// ListAllFiles returns all files by request, even if it requires several API
// calls.
func (client *Client) ListAllFiles(
	projectID string,
	request FilesListRequest,
) ([]File, error) {
	result := []File{}

	for {
		files, err := client.ListFiles(projectID, request)
		if err != nil {
			return nil, err
		}

		result = append(result, files.Items...)

		if request.Cursor.Limit > 0 {
			request.Cursor.Limit -= len(files.Items)

			if request.Cursor.Limit == 0 {
				return result, nil
			}
		}

		if request.Cursor.Offset+len(files.Items) < files.TotalCount {
			request.Cursor.Offset += len(files.Items)
		} else {
			break
		}

		client.Logger.Infof(
			"<= %d/%d files retrieved",
			request.Cursor.Offset,
			files.TotalCount,
		)
	}

	client.Logger.Infof(
		"<= %d/%d files retrieved",
		len(result),
		len(result),
	)

	return result, nil
}
