package smartling

import (
	"fmt"
	"net/url"
)

// FilesListRequest represents request used to filter files returned by
// list files API call.
type FilesListRequest struct {
	// Cursor is a limit/offset pair, used to paginate reply.
	Cursor LimitOffsetRequest

	// URIMask instructs API to return only files with a URI containing the
	// given substring. Case is ignored.
	URIMask string

	// FileTypes instructs API to return only specified file types.
	FileTypes []FileType

	// LastUploadedAfter instructs API to return files uploaded after specified
	// date.
	LastUploadedAfter UTC

	// LastUploadedBefore instructs API to return files uploaded after
	// specified date.
	LastUploadedBefore UTC
}

// GetQuery returns URL values representation of files list request.
func (request *FilesListRequest) GetQuery() url.Values {
	query := request.Cursor.GetQuery()
	if len(request.URIMask) > 0 {
		query.Set("uriMask", request.URIMask)
	}

	for _, fileType := range request.FileTypes {
		query.Add("fileTypes[]", fmt.Sprint(fileType))
	}

	if !request.LastUploadedAfter.IsZero() {
		query.Set("lastUploadedAfter", request.LastUploadedAfter.String())
	}

	if !request.LastUploadedBefore.IsZero() {
		query.Set("lastUploadedBefore", request.LastUploadedBefore.String())
	}

	return query
}
