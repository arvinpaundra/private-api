package constant

import "errors"

var (
	ErrSubmissionNotFound       = errors.New("submission not found")
	ErrSubmissionAnswerNotFound = errors.New("submission answer not found")

	ErrInvalidStatus         = errors.New("invalid submission status")
	ErrCannotSubmit          = errors.New("cannot submit submission in current state")
	ErrCannotCancel          = errors.New("cannot cancel submission in current state")
	ErrSubmissionAlreadyDone = errors.New("submission already submitted or canceled")
	ErrDuplicateAnswer       = errors.New("answer for this question already submitted")

	// Context mapping errors - submission's perspective on related entities
	ErrModuleNotFound   = errors.New("module not found")
	ErrQuestionNotFound = errors.New("question not found")
	ErrChoiceNotFound   = errors.New("choice not found")
)
