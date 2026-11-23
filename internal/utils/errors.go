package utils

type ClientError struct {
	Message string
}

// Error() untuk wrap error di service agar bisa terbaca ke handler
func (e *ClientError) Error() string {
	return e.Message
}

func NewClientError(msg string) error {
	return &ClientError{
		Message: msg,
	}
}
