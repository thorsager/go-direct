package godirect

import "fmt"

type NotFoundError struct {
	path string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("redirect not found for: %s", e.path)
}

func NotFound(path string) *NotFoundError {
	return &NotFoundError{path: path}
}

type NotNumberError struct {
	str string
}

func (e *NotNumberError) Error() string {
	return fmt.Sprintf("Not a base64 encoded number: %s", e.str)
}

func IsNotFound(err error) bool {
	_, ok := err.(*NotFoundError)
	return ok
}
