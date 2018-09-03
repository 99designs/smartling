package smartling

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type APIError struct {
	Cause    error
	Code     string
	URL      string
	Params   url.Values
	Payload  []byte
	Response []byte
	Headers  *http.Header
}

func (err APIError) Error() string {
	url := err.URL
	if len(err.Params) > 0 {
		url += "?" + err.Params.Encode()
	}

	code := ""
	if err.Code != "" {
		code = fmt.Sprintf(" [code %s]", code)
	}

	headers := &bytes.Buffer{}
	_ = err.Headers.Write(headers)

	return fmt.Sprintf(
		"%s%s\nURL: %s%s%s%s",
		err.Cause,
		code,
		url,
		err.section("Request Body", string(err.Payload)),
		err.section("Response Body", string(err.Response)),
		err.section("Response Headers", headers.String()),
	)
}

func (err *APIError) section(header string, contents string) string {
	if contents == "" {
		return ""
	}

	return fmt.Sprintf(
		"\n\n%s:\n%s",
		header,
		regexp.MustCompile(`(?m)^`).ReplaceAllLiteralString(
			strings.TrimSuffix(contents, "\n"),
			"  ",
		),
	)
}
