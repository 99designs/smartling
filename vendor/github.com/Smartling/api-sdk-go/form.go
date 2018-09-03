package smartling

import (
	"bytes"
	"mime/multipart"
)

type Form struct {
	Writer *multipart.Writer
	Body   *bytes.Buffer
}

func (form *Form) Close() error {
	return form.Writer.Close()
}

func (form *Form) GetContentType() string {
	return form.Writer.FormDataContentType()
}

func (form *Form) Bytes() []byte {
	return form.Body.Bytes()
}
