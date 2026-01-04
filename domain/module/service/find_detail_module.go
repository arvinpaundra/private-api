package service

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/module/repository"
	"github.com/arvinpaundra/private-api/domain/module/response"
	"github.com/arvinpaundra/private-api/domain/shared/interfaces"
)

type FindDetailModuleCommand struct {
	Slug string `form:"slug"`
}

type FindDetailModule struct {
	authStorage  interfaces.AuthenticatedUser
	moduleReader repository.ModuleReader
}

func NewFindDetailModule(
	authStorage interfaces.AuthenticatedUser,
	moduleReader repository.ModuleReader,
) *FindDetailModule {
	return &FindDetailModule{
		authStorage:  authStorage,
		moduleReader: moduleReader,
	}
}

func (s *FindDetailModule) Execute(ctx context.Context, command *FindDetailModuleCommand) (*response.Module, error) {
	module, err := s.moduleReader.FindBySlug(ctx, command.Slug, s.authStorage.GetUserId())
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
