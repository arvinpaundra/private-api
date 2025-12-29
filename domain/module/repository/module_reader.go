package repository

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/module/entity"
)

type ModuleReader interface {
	FindByID(ctx context.Context, moduleID, userID string) (*entity.Module, error)
	FindBySlug(ctx context.Context, slug, userID string) (*entity.Module, error)
	FindModuleDetailBySlug(ctx context.Context, slug, userID string) (*entity.Module, error)
	TotalModules(ctx context.Context, userID, subjectID, gradeID, keyword string) (int, error)
	FindAllModules(ctx context.Context, userID, subjectID, gradeID, keyword string, limit, offset int) ([]*entity.Module, error)
}
