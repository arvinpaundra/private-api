package service

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/module/repository"
	"github.com/arvinpaundra/private-api/domain/module/response"
)

type FindPublishedModuleCommand struct {
	Slug string `form:"slug"`
}

type FindPublishedModule struct {
	moduleReader repository.ModuleReader
}

func NewFindPublishedModule(
	moduleReader repository.ModuleReader,
) *FindPublishedModule {
	return &FindPublishedModule{
		moduleReader: moduleReader,
	}
}

func (s *FindPublishedModule) Execute(ctx context.Context, command *FindPublishedModuleCommand) (*response.Module, error) {
	module, err := s.moduleReader.FindPublishedModuleBySlug(ctx, command.Slug)
	if err != nil {
		return nil, err
	}

	totalQuestions, err := s.moduleReader.CountQuestionsByModuleSlug(ctx, module.Slug)
	if err != nil {
		return nil, err
	}

	result := &response.Module{
		ID:             module.ID,
		UserID:         module.UserID,
		SubjectID:      module.SubjectID,
		GradeID:        module.GradeID,
		Slug:           module.Slug,
		Title:          module.Title,
		Type:           module.Type,
		IsPublished:    module.IsPublished,
		QuestionsCount: totalQuestions,
	}

	return result, nil
}
