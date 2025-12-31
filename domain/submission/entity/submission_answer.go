package entity

import (
	"github.com/arvinpaundra/private-api/core/trait"
	"github.com/arvinpaundra/private-api/core/util"
)

type SubmissionAnswer struct {
	trait.Createable
	trait.Updateable
	trait.Removeable

	ID           string
	SubmissionID string
	QuestionSlug string
	Question     string
	Answer       string
	IsCorrect    bool
}

func NewSubmissionAnswer(submissionID, questionSlug, question, answer string, isCorrect bool) *SubmissionAnswer {
	submissionAnswer := &SubmissionAnswer{
		ID:           util.GenerateUUID(),
		SubmissionID: submissionID,
		QuestionSlug: questionSlug,
		Question:     question,
		Answer:       answer,
		IsCorrect:    isCorrect,
	}

	submissionAnswer.MarkCreate()

	return submissionAnswer
}

func (sa *SubmissionAnswer) UpdateCorrectness(isCorrect bool) {
	sa.IsCorrect = isCorrect
	sa.MarkUpdate()
}
