package smartling

type FileUploadRequest struct {
	FileURIRequest

	File      []byte
	FileType  FileType
	Authorize bool

	LocalesToAuthorize []string

	Smartling struct {
		Namespace   string
		FileCharset string
		Directives  map[string]string
	}
}

func (request *FileUploadRequest) GetForm() (*Form, error) {
	form, err := request.FileURIRequest.GetForm()
	if err != nil {
		return nil, err
	}

	writer, err := form.Writer.CreateFormFile("file", request.FileURI)
	if err != nil {
		return nil, err
	}

	_, err = writer.Write(request.File)
	if err != nil {
		return nil, err
	}

	err = form.Writer.WriteField("fileType", string(request.FileType))
	if err != nil {
		return nil, err
	}

	if request.Authorize {
		err = form.Writer.WriteField("authorize", "true")
		if err != nil {
			return nil, err
		}
	}

	for _, locale := range request.LocalesToAuthorize {
		err = form.Writer.WriteField("localeIdsToAuthorize[]", locale)
		if err != nil {
			return nil, err
		}
	}

	if request.Smartling.Namespace != "" {
		err = form.Writer.WriteField(
			"smartling.namespace",
			request.Smartling.Namespace,
		)
		if err != nil {
			return nil, err
		}
	}

	if request.Smartling.FileCharset != "" {
		err = form.Writer.WriteField(
			"smartling.file_charset",
			request.Smartling.FileCharset,
		)
		if err != nil {
			return nil, err
		}
	}

	for directive, value := range request.Smartling.Directives {
		err = form.Writer.WriteField("smartling."+directive, value)
		if err != nil {
			return nil, err
		}
	}

	return form, nil
}
