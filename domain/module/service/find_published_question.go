package service

import (
	"context"
	"errors"

	"github.com/arvinpaundra/private-api/domain/module/constant"
	"github.com/arvinpaundra/private-api/domain/module/repository"
	"github.com/arvinpaundra/private-api/domain/module/response"
)

type FindPublishedQuestionCommand struct {
	ModuleSlug   string `form:"-" validate:"required"`
	QuestionSlug string `form:"-" validate:"required"`
}

type FindPublishedQuestion struct {
	moduleReader repository.ModuleReader
}

func NewFindPublishedQuestion(
	moduleReader repository.ModuleReader,
) *FindPublishedQuestion {
	return &FindPublishedQuestion{
		moduleReader: moduleReader,
	}
}

func (s *FindPublishedQuestion) Execute(ctx context.Context, command *FindPublishedQuestionCommand) (*response.QuestionDetail, error) {
	// Verify module exists and is published
	module, err := s.moduleReader.FindPublishedModuleBySlug(ctx, command.ModuleSlug)
	if err != nil {
		return nil, err
	}

	// Find the specific question
	question, err := s.moduleReader.FindPublishedQuestionBySlug(ctx, module.Slug, command.QuestionSlug)
	if err != nil {
		return nil, err
	}

	// Get next question for cursor
	nextQuestion, err := s.moduleReader.FindNextPublishedQuestion(ctx, module.Slug, command.QuestionSlug)
	if err != nil && !errors.Is(err, constant.ErrQuestionNotFound) {
		return nil, err
	}

	var nextQuestionSlug *string

	if nextQuestion != nil {
		nextQuestionSlug = &nextQuestion.Slug
	}

	// Build response
	choices := make([]*response.Choice, len(question.Choices))

	for i, choice := range question.Choices {
		choices[i] = &response.Choice{
			ID:      choice.ID,
			Content: choice.Content,
		}
	}

	return &response.QuestionDetail{
		ID:               question.ID,
		Content:          question.Content,
		Slug:             question.Slug,
		Choices:          choices,
		NextQuestionSlug: nextQuestionSlug,
	}, nil
}
