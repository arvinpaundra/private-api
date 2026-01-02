package repository

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/submission/entity"
)

type ModuleACL interface {
	GetCorrectAnswer(ctx context.Context, moduleSlug, questionSlug string) (*entity.Choice, error)
	GetNextQuestionSlug(ctx context.Context, moduleSlug, currentQuestionSlug string) (*string, error)
	GetPublishedModule(ctx context.Context, moduleSlug string) (*entity.Module, error)
	GetQuestionBySlug(ctx context.Context, moduleSlug, questionSlug string) (*entity.Question, error)
	GetTotalQuestions(ctx context.Context, moduleSlug string) (int, error)
	GetAllPublishedModules(ctx context.Context, keyword string) ([]*entity.Module, error)
}
