package client

type UnauthorizedError struct{}

func (UnauthorizedError) Error() string {
	return "Item is not accessible"
}

func IsUnauthorizedError(err error) bool {
	_, ok := err.(UnauthorizedError)
	return ok
}

type RateLimitError struct{}

func (RateLimitError) Error() string {
	return "Over rate limit"
}

func IsRateLimitError(err error) bool {
	_, ok := err.(RateLimitError)
	return ok
}

type DeletedError struct{}

func (DeletedError) Error() string {
	return "Item is deleted"
}

func IsDeletedError(err error) bool {
	_, ok := err.(DeletedError)
	return ok
}
