package errorsx

import "errors"

// retryableError represents an error with a retryable flag
type retryableError string

// NewRetryableError creates a retryable tool error
func NewRetryableError(text string) error {
	return retryableError(text)
}

// Error returns the error text
func (e retryableError) Error() string {
	return string(e)
}

// IsRetryable checks if an error is retryable
func IsRetryable(err error) bool {
	var e *retryableError
	return errors.As(err, &e)
}
