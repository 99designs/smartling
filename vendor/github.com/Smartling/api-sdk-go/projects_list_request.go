package smartling

import (
	"fmt"
	"net/url"
)

// ProjectsListRequest is a request used in GetProjectsList method.
type ProjectsListRequest struct {
	// Cursor specifies limit/offset pagination pair.
	Cursor LimitOffsetRequest

	// ProjectNameFilter specifies filter for project name.
	ProjectNameFilter string

	// IncludeArchived specifies should archived items be included or not.
	IncludeArchived bool
}

// GetQuery returns URL-encoded representation of current request.
func (request ProjectsListRequest) GetQuery() url.Values {
	query := request.Cursor.GetQuery()

	if len(request.ProjectNameFilter) > 0 {
		query.Set("projectNameFilter", request.ProjectNameFilter)
	}

	query.Set("includeArchived", fmt.Sprint(request.IncludeArchived))

	return query
}
