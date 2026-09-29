package utils

type ClientError struct {
	Message string
}

func (e *ClientError) Error() string {
	return e.Message
}

func NewClientError(msg string) error {
	return &ClientError{
		Message: msg,
	}
}