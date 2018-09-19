package smartling

import (
	"fmt"
	"net/url"
)

// RetrievalType describes type of file download.
// https://help.smartling.com/v1.0/reference#get_projects-projectid-locales-localeid-file
type RetrievalType string

const (
	// RetrieveDefault specifies that Smartling will decide what type of
	// translation will be returned.
	RetrieveDefault RetrievalType = ""

	// RetrievePending specifies that Smartling returns any translations
	// (including non-published translations)
	RetrievePending = "pending"

	// RetrievePublished specifies that Smartling returns only
	// published/pre-published translations.
	RetrievePublished = "published"

	// RetrievePseudo specifies that Smartling returns a modified version of
	// the original text with certain characters transformed and the text
	// expanded.
	RetrievePseudo = "pseudo"

	// RetrieveChromeInstrumented specifies that Smartling returns a modified
	// version of the original file with strings wrapped in a specific set of
	// Unicode symbols that can later be recognized and matched by the Chrome
	// Context Capture Extension
	RetrieveChromeInstrumented = "contextMatchingInstrumented"
)

// FileDownloadRequest represents optional parameters for file download
// operation.
type FileDownloadRequest struct {
	FileURIRequest

	Type            RetrievalType
	IncludeOriginal bool
}

// GetQuery returns URL values representation of download file query params.
func (request FileDownloadRequest) GetQuery() url.Values {
	query := request.FileURIRequest.GetQuery()

	query.Set("fileUri", request.FileURI)

	if request.Type != RetrieveDefault {
		query.Set("retrievalType", fmt.Sprint(request.Type))
	}

	if request.IncludeOriginal {
		query.Set("includeOriginalStrings", "true")
	}

	return query
}
