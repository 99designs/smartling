package smartling

import (
	"fmt"
	"net/url"
)

// LimitOffsetRequest is a base request for all other requests to set
// pagination options,  e.g. limit and offset.
type LimitOffsetRequest struct {
	Offset int
	Limit  int
}

// GetQuery returns URL-encoded representation of current request.
func (request LimitOffsetRequest) GetQuery() url.Values {
	query := url.Values{}

	query.Set("offset", fmt.Sprint(request.Offset))

	if request.Limit > 0 {
		query.Set("limit", fmt.Sprint(request.Limit))
	}

	return query
}
