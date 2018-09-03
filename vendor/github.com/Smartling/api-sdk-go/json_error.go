package smartling

import "fmt"

type JSONError struct {
	Cause    error
	Response []byte
}

func (err JSONError) Error() string {
	return fmt.Sprintf(
		"unable to parse reply as JSON: %s\n%s",
		err.Cause,
		err.Response,
	)
}
