package service

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/module/repository"
	"github.com/arvinpaundra/private-api/domain/shared/interfaces"
)

type TogglePublishModuleCommand struct {
	Slug string `json:"-" validate:"required"`
}

type TogglePublishModule struct {
	authStorage  interfaces.AuthenticatedUser
	moduleReader repository.ModuleReader
	moduleWriter repository.ModuleWriter
}

func NewTogglePublishModule(
	authStorage interfaces.AuthenticatedUser,
	moduleReader repository.ModuleReader,
	moduleWriter repository.ModuleWriter,
) *TogglePublishModule {
	return &TogglePublishModule{
		authStorage:  authStorage,
		moduleReader: moduleReader,
		moduleWriter: moduleWriter,
	}
}

func (s *TogglePublishModule) Execute(ctx context.Context, command *TogglePublishModuleCommand) error {
	// Find module by slug
	module, err := s.moduleReader.FindBySlug(ctx, command.Slug, s.authStorage.GetUserId())
	if err != nil {
		return err
	}

	// Toggle publish state
	if module.IsPublished {
		module.Unpublish()
	} else {
		module.Publish()
	}

	// Save via aggregate root guard
	if err := s.moduleWriter.Save(ctx, module); err != nil {
		return err
	}

	return nil
}
