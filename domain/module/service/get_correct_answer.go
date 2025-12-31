package service

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/module/repository"
	"github.com/arvinpaundra/private-api/domain/module/response"
)

type GetCorrectAnswerCommand struct {
	ModuleSlug   string `validate:"required"`
	QuestionSlug string `validate:"required"`
}

type GetCorrectAnswer struct {
	moduleReader repository.ModuleReader
}

func NewGetCorrectAnswer(
	moduleReader repository.ModuleReader,
) *GetCorrectAnswer {
	return &GetCorrectAnswer{
		moduleReader: moduleReader,
	}
}

func (s *GetCorrectAnswer) Execute(ctx context.Context, command *GetCorrectAnswerCommand) (*response.ChoiceWithAnswer, error) {
	// Verify module exists and is published
	_, err := s.moduleReader.FindPublishedModuleBySlug(ctx, command.ModuleSlug)
	if err != nil {
		return nil, err
	}

	// Find the specific question
	question, err := s.moduleReader.FindPublishedQuestionBySlug(ctx, command.ModuleSlug, command.QuestionSlug)
	if err != nil {
		return nil, err
	}

	// Find the correct answer from choices
	for _, choice := range question.Choices {
		if choice.IsCorrectAnswer {
			return &response.ChoiceWithAnswer{
				ID:              choice.ID,
				Content:         choice.Content,
				IsCorrectAnswer: choice.IsCorrectAnswer,
			}, nil
		}
	}

	// No correct answer found
	return nil, nil
}
