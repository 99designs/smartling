package smartling

// RenameFileRequest represents fileUri query parameter, commonly used in API.
type RenameFileRequest struct {
	FileURIRequest

	NewFileURI string
}

func (request *RenameFileRequest) GetForm() (*Form, error) {
	form, err := request.FileURIRequest.GetForm()
	if err != nil {
		return nil, err
	}

	err = form.Writer.WriteField("newFileUri", request.NewFileURI)
	if err != nil {
		return nil, err
	}

	return form, nil
}
