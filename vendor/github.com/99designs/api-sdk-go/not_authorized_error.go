package smartling

type NotAuthorizedError struct{}

func (err NotAuthorizedError) Error() string {
	return "authentication parameters are invalid"
}
