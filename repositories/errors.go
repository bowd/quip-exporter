package repositories

type NotFoundError struct{}

func (NotFoundError) Error() string {
	return "Item not found in repository"
}

func IsNotFoundError(err error) bool {
	_, ok := err.(NotFoundError)
	return ok
}
