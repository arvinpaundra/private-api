package constant

import "errors"

var (
	ErrGradeAlreadyExists = errors.New("grade with the same name already exists")
	ErrGradeNotFound      = errors.New("grade not found")
)
