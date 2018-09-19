package smartling

import (
	"bytes"
	"mime/multipart"
	"net/url"
)

// FileURIRequest represents fileUri query parameter, commonly used in API.
type FileURIRequest struct {
	FileURI string
}

// GetQuery returns URL value representation for file URI.
func (request FileURIRequest) GetQuery() url.Values {
	query := url.Values{}

	query.Set("fileUri", request.FileURI)

	return query
}

func (request *FileURIRequest) GetForm() (*Form, error) {
	var (
		body   = &bytes.Buffer{}
		writer = multipart.NewWriter(body)
	)

	err := writer.WriteField("fileUri", request.FileURI)
	if err != nil {
		return nil, err
	}

	return &Form{
		Writer: writer,
		Body:   body,
	}, nil
}
