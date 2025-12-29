package entity

import (
	"github.com/arvinpaundra/private-api/core/trait"
	"github.com/arvinpaundra/private-api/core/util"
	"github.com/arvinpaundra/private-api/domain/module/constant"
)

type Question struct {
	trait.Createable
	trait.Updateable
	trait.Removeable

	ID       string
	ModuleID string
	Content  string
	Slug     string

	Choices []*QuestionChoice
}

func NewQuestion(moduleID, content string) *Question {
	question := &Question{
		ID:       util.GenerateUUID(),
		ModuleID: moduleID,
		Content:  content,
	}

	question.MarkCreate()

	return question
}

func (q *Question) GenSlug() error {
	slug, err := util.RandomAlphanumeric(12)
	if err != nil {
		return err
	}

	q.Slug = slug

	return nil
}

func (q *Question) IsValidChoices() error {
	if len(q.Choices) < 2 {
		return constant.ErrMinTwoChoices
	} else if len(q.Choices) > 4 {
		return constant.ErrMaxFourChoices
	}

	hasCorrectAnswer := false

	for _, choice := range q.Choices {
		if choice.IsCorrectAnswer {
			// only one correct answer is allowed
			if hasCorrectAnswer {
				return constant.ErrMultipleCorrectAnswers
			}

			hasCorrectAnswer = true
		}
	}

	if !hasCorrectAnswer {
		return constant.ErrNoCorrectAnswer
	}

	return nil
}

func (q *Question) AddChoice(choice *QuestionChoice) {
	q.Choices = append(q.Choices, choice)
}

type QuestionChoice struct {
	trait.Createable
	trait.Updateable
	trait.Removeable

	ID              string
	QuestionID      string
	Content         string
	IsCorrectAnswer bool
}

func NewQuestionChoice(questionID, content string) *QuestionChoice {
	choice := &QuestionChoice{
		ID:         util.GenerateUUID(),
		QuestionID: questionID,
		Content:    content,
	}

	choice.MarkCreate()

	return choice
}

func (qc *QuestionChoice) SetAsCorrectAnswer() {
	qc.IsCorrectAnswer = true
}
