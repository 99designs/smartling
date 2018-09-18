package smartling

import (
	"fmt"
)

const (
	endpointFileRename = "/files-api/v2/projects/%s/file/rename"
)

// RenameFile renames file to new URI.
func (client *Client) RenameFile(
	projectID string,
	oldURI string,
	newURI string,
) error {
	request := RenameFileRequest{}
	request.FileURI = oldURI
	request.NewFileURI = newURI

	form, err := request.GetForm()
	if err != nil {
		return fmt.Errorf(
			"failed to create file rename form: %s",
			err,
		)
	}

	err = form.Close()
	if err != nil {
		return fmt.Errorf(
			"unable to close file rename form: %s", err,
		)
	}

	_, _, err = client.Post(
		fmt.Sprintf(endpointFileRename, projectID),
		form.Bytes(),
		nil,
		ContentTypeOption(form.GetContentType()),
	)
	if err != nil {
		return fmt.Errorf(
			"failed to rename file: %s", err,
		)
	}

	return nil
}
