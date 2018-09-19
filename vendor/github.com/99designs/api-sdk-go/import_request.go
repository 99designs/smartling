package smartling

type TranslationState string

const (
	TranslationStatePublished       TranslationState = "PUBLISHED"
	TranslationStatePostTranslation                  = "POST_TRANSLATION"
)

type ImportRequest struct {
	FileURIRequest

	File             []byte
	FileType         FileType
	TranslationState TranslationState
	Overwrite        bool
}

func (request *ImportRequest) GetForm() (*Form, error) {
	form, err := request.FileURIRequest.GetForm()
	if err != nil {
		return nil, err
	}

	file, err := form.Writer.CreateFormFile("file", request.FileURI)
	if err != nil {
		return nil, err
	}

	_, err = file.Write(request.File)
	if err != nil {
		return nil, err
	}

	err = form.Writer.WriteField("fileType", string(request.FileType))
	if err != nil {
		return nil, err
	}

	err = form.Writer.WriteField(
		"translationState",
		string(request.TranslationState),
	)
	if err != nil {
		return nil, err
	}

	if request.Overwrite {
		err = form.Writer.WriteField("overwrite", "true")
		if err != nil {
			return nil, err
		}
	}

	return form, nil
}
