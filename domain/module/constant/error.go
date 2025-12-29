package constant

import "errors"

var (
	ErrModuleNotFound   = errors.New("module not found")
	ErrQuestionNotFound = errors.New("question not found")
	ErrChoiceNotFound   = errors.New("choice not found")

	ErrMinTwoChoices          = errors.New("a question must have at least two choices")
	ErrMaxFourChoices         = errors.New("a question must not have more than four choices")
	ErrMultipleCorrectAnswers = errors.New("a question must not have more than one correct answer")
	ErrNoCorrectAnswer        = errors.New("a question must have at least one correct answer")

	// Context mapping errors - module's perspective on related entities
	ErrSubjectNotFound = errors.New("subject not found")
	ErrGradeNotFound   = errors.New("grade not found")
)
