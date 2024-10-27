package app_errors

type NotFoundError struct {
	Message string
}

func (receiver NotFoundError) Error() string {
	return receiver.Message
}

func NewNotFoundError(message string) NotFoundError {
	if message == "" {
		message = "Resource not found"
	}
	return NotFoundError{
		Message: message,
	}
}
