package smartling

import (
	"fmt"
)

const (
	endpointUploadFile = "/files-api/v2/projects/%s/file"
)

type FileUploadResult struct {
	Overwritten bool
	StringCount int
	WordCount   int
}

// DownloadFile downloads original file from project.
func (client *Client) UploadFile(
	projectID string,
	request FileUploadRequest,
) (*FileUploadResult, error) {
	var result FileUploadResult

	form, err := request.GetForm()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create file upload form: %s",
			err,
		)
	}

	err = form.Close()
	if err != nil {
		return nil, fmt.Errorf(
			"failed to close upload file form: %s", err,
		)
	}

	_, _, err = client.Post(
		fmt.Sprintf(endpointUploadFile, projectID),
		form.Bytes(),
		&result,
		ContentTypeOption(form.GetContentType()),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to download original file: %s",
			err,
		)
	}

	return &result, nil
}
