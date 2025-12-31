package service

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/module/repository"
	"github.com/arvinpaundra/private-api/domain/module/response"
)

type ValidatePublishedModuleCommand struct {
	ModuleSlug string `validate:"required"`
}

type ValidatePublishedModule struct {
	moduleReader repository.ModuleReader
}

func NewValidatePublishedModule(
	moduleReader repository.ModuleReader,
) *ValidatePublishedModule {
	return &ValidatePublishedModule{
		moduleReader: moduleReader,
	}
}

func (s *ValidatePublishedModule) Execute(ctx context.Context, command *ValidatePublishedModuleCommand) (*response.Module, error) {
	// Verify module exists and is published
	module, err := s.moduleReader.FindPublishedModuleBySlug(ctx, command.ModuleSlug)
	if err != nil {
		return nil, err
	}

	return &response.Module{
		ID:          module.ID,
		UserID:      module.UserID,
		SubjectID:   module.SubjectID,
		GradeID:     module.GradeID,
		Title:       module.Title,
		Description: module.Description,
		Type:        module.Type,
		IsPublished: module.IsPublished,
		Slug:        module.Slug,
	}, nil
}
