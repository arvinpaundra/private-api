package service

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/module/repository"
	"github.com/arvinpaundra/private-api/domain/shared/interfaces"
)

type DeleteModuleCommand struct {
	Slug string `json:"-" validate:"required"`
}

type DeleteModule struct {
	authStorage  interfaces.AuthenticatedUser
	moduleReader repository.ModuleReader
	moduleWriter repository.ModuleWriter
}

func NewDeleteModule(
	authStorage interfaces.AuthenticatedUser,
	moduleReader repository.ModuleReader,
	moduleWriter repository.ModuleWriter,
) *DeleteModule {
	return &DeleteModule{
		authStorage:  authStorage,
		moduleReader: moduleReader,
		moduleWriter: moduleWriter,
	}
}

func (s *DeleteModule) Execute(ctx context.Context, command *DeleteModuleCommand) error {
	// Find module by slug
	module, err := s.moduleReader.FindBySlug(ctx, command.Slug, s.authStorage.GetUserId())
	if err != nil {
		return err
	}

	// Mark module as Deleted
	module.MarkRemove()

	// Save via aggregate root guard
	if err := s.moduleWriter.Save(ctx, module); err != nil {
		return err
	}

	return nil
}
