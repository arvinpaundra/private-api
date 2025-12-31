package service

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/module/repository"
)

type CountModuleQuestionsCommand struct {
	ModuleSlug string `validate:"required"`
}

type CountModuleQuestions struct {
	moduleReader repository.ModuleReader
}

func NewCountModuleQuestions(
	moduleReader repository.ModuleReader,
) *CountModuleQuestions {
	return &CountModuleQuestions{
		moduleReader: moduleReader,
	}
}

func (s *CountModuleQuestions) Execute(ctx context.Context, command *CountModuleQuestionsCommand) (int, error) {
	// Get total questions count for the module
	total, err := s.moduleReader.CountQuestionsByModuleSlug(ctx, command.ModuleSlug)
	if err != nil {
		return 0, err
	}

	return total, nil
}
