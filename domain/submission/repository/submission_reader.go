package repository

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/submission/entity"
)

type SubmissionReader interface {
	FindByCode(ctx context.Context, code string) (*entity.Submission, error)
	CountSubmitted(ctx context.Context) (int, error)
	TotalSubmissions(ctx context.Context, moduleID, status, keyword string) (int, error)
	FindAllSubmissions(ctx context.Context, moduleID, status, keyword string, limit, offset int) ([]*entity.Submission, error)
	FindAllSubmittedGroupedByModule(ctx context.Context, moduleIDs []string) (map[string][]*entity.Submission, error)
}
