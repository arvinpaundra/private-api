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

func NewQuestion(moduleID, content string) (*Question, error) {
	question := &Question{
		ID:       util.GenerateUUID(),
		ModuleID: moduleID,
		Content:  content,
	}

	err := question.GenSlug()
	if err != nil {
		return nil, err
	}

	question.MarkCreate()

	return question, nil
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
	counter := 0

	hasCorrectAnswer := false

	for _, choice := range q.Choices {
		if choice.IsRemoved() {
			continue
		}

		if choice.IsCorrectAnswer {
			// only one correct answer is allowed
			if hasCorrectAnswer {
				return constant.ErrMultipleCorrectAnswers
			}

			hasCorrectAnswer = true
		}

		counter++
	}

	if !hasCorrectAnswer {
		return constant.ErrNoCorrectAnswer
	}

	if counter < 2 {
		return constant.ErrMinTwoChoices
	} else if counter > 4 {
		return constant.ErrMaxFourChoices
	}

	return nil
}

func (q *Question) AddChoice(choice *QuestionChoice) {
	q.Choices = append(q.Choices, choice)
}

func (q *Question) UpdateContent(content string) {
	q.Content = content
	q.MarkUpdate()
}

func (q *Question) ClearChoices() {
	// Mark existing choices for removal before clearing
	for _, choice := range q.Choices {
		choice.MarkRemove()
	}

	q.MarkUpdate()
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
