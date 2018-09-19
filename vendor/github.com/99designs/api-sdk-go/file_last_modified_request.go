package smartling

type FileLastModifiedRequest struct {
	FileURIRequest

	LastModifiedAfter UTC
}

func (request *FileLastModifiedRequest) GetForm() (*Form, error) {
	form, err := request.FileURIRequest.GetForm()

	if err != nil {
		return nil, err
	}

	err = form.Writer.WriteField("lastModifiedAfter", request.LastModifiedAfter.String())

	if err != nil {
		return nil, err
	}

	return form, nil
}
