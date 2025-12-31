package repository

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/module/entity"
)

type ModuleReader interface {
	FindByID(ctx context.Context, moduleID, userID string) (*entity.Module, error)
	FindBySlug(ctx context.Context, slug, userID string) (*entity.Module, error)
	FindModuleDetailBySlug(ctx context.Context, slug, userID string) (*entity.Module, error)
	FindPublishedModuleBySlug(ctx context.Context, slug string) (*entity.Module, error)
	FindPublishedQuestionBySlug(ctx context.Context, moduleSlug, questionSlug string) (*entity.Question, error)
	FindNextPublishedQuestion(ctx context.Context, moduleSlug, currentQuestionSlug string) (*entity.Question, error)
	CountQuestionsByModuleSlug(ctx context.Context, moduleSlug string) (int, error)
	TotalModules(ctx context.Context, userID, subjectID, gradeID, keyword string) (int, error)
	FindAllModules(ctx context.Context, userID, subjectID, gradeID, keyword string, limit, offset int) ([]*entity.Module, error)
}
