package smartling

import (
	"fmt"
)

const (
	endpointFileDelete = "/files-api/v2/projects/%s/file/delete"
)

// DeleteFile removes specified files from project.
func (client *Client) DeleteFile(
	projectID string,
	uri string,
) error {
	request := FileURIRequest{
		FileURI: uri,
	}

	form, err := request.GetForm()
	if err != nil {
		return fmt.Errorf(
			"failed to create file delete form: %s",
			err,
		)
	}

	err = form.Close()
	if err != nil {
		return fmt.Errorf(
			"failed to close file delete form: %s", err,
		)
	}

	_, _, err = client.Post(
		fmt.Sprintf(endpointFileDelete, projectID),
		form.Bytes(),
		nil,
		ContentTypeOption(form.GetContentType()),
	)
	if err != nil {
		return fmt.Errorf(
			"failed to remove file: %s", err,
		)
	}

	return nil
}
