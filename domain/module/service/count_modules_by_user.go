package service

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/module/repository"
)

type CountModulesByUser struct {
	moduleReader repository.ModuleReader
}

func NewCountModulesByUser(moduleReader repository.ModuleReader) *CountModulesByUser {
	return &CountModulesByUser{
		moduleReader: moduleReader,
	}
}

func (s *CountModulesByUser) Execute(ctx context.Context, userID string) (int, error) {
	count, err := s.moduleReader.CountByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}

	return count, nil
}
