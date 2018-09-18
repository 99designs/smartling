package smartling

type NotFoundError struct{}

func (err NotFoundError) Error() string {
	return "requested entity is not found"
}
