package entity

type ValidationError struct {
	innerError error
}

func NewValidationError(err error) *ValidationError {
	return &ValidationError{
		innerError: err,
	}
}

func (e *ValidationError) Error() string {
	return e.innerError.Error()
}

func (e *ValidationError) Unwrap() error {
	return e.innerError
}
