package service

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/module/entity"
	"github.com/arvinpaundra/private-api/domain/module/repository"
	"github.com/arvinpaundra/private-api/domain/shared/interfaces"
)

type AddQuestionsCommand struct {
	ModuleSlug string         `json:"-" validate:"required"`
	Questions  []*AddQuestion `json:"questions" validate:"required,min=1,dive"`
}

type AddQuestion struct {
	Content string               `json:"content" validate:"required"`
	Choices []*AddQuestionChoice `json:"choices" validate:"required,min=2,max=4,dive"`
}

type AddQuestionChoice struct {
	Content         string `json:"content" validate:"required"`
	IsCorrectAnswer bool   `json:"is_correct_answer"`
}

type AddQuestions struct {
	authStorage  interfaces.AuthenticatedUser
	moduleReader repository.ModuleReader
	uow          repository.UnitOfWork
}

func NewAddQuestions(
	authStorage interfaces.AuthenticatedUser,
	moduleReader repository.ModuleReader,
	uow repository.UnitOfWork,
) *AddQuestions {
	return &AddQuestions{
		authStorage:  authStorage,
		moduleReader: moduleReader,
		uow:          uow,
	}
}

func (s *AddQuestions) Execute(ctx context.Context, command *AddQuestionsCommand) error {
	// Check if module exists and belongs to the user
	module, err := s.moduleReader.FindBySlug(ctx, command.ModuleSlug, s.authStorage.GetUserId())
	if err != nil {
		return err
	}

	// Process each question
	for _, questionCmd := range command.Questions {
		// Create new question
		question := entity.NewQuestion(module.ID, questionCmd.Content)

		// Generate slug for question
		err = question.GenSlug()
		if err != nil {
			return err
		}

		// Add choices to question
		for _, choiceCmd := range questionCmd.Choices {
			choice := entity.NewQuestionChoice(question.ID, choiceCmd.Content)
			if choiceCmd.IsCorrectAnswer {
				choice.SetAsCorrectAnswer()
			}

			question.AddChoice(choice)
		}

		// Validate choices
		err = question.IsValidChoices()
		if err != nil {
			return err
		}

		// Add question to module
		module.AddQuestion(question)
	}

	// Begin transaction
	tx, err := s.uow.Begin()
	if err != nil {
		return err
	}

	// Save module (which will cascade to questions)
	if err := tx.ModuleWriter().Save(ctx, module); err != nil {
		if uowErr := tx.Rollback(); uowErr != nil {
			return uowErr
		}
		return err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
