package constant

import "errors"

var (
	ErrSubjectAlreadyExists = errors.New("subject with the same name already exists")
	ErrSubjectNotFound      = errors.New("subject not found")
)
