package repository

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/subject/entity"
)

type SubjectReader interface {
	HasSimilarSubject(ctx context.Context, name string, userID string) (bool, error)
	HasSimilarSubjectExclusive(ctx context.Context, name string, userID string, excludeSubjectID string) (bool, error)
	IsSubjectExist(ctx context.Context, subjectID string, userID string) (bool, error)
	FindSubjectByID(ctx context.Context, subjectID string, userID string) (*entity.Subject, error)
	CountByUserID(ctx context.Context, userID string) (int, error)
	AllSubjects(ctx context.Context, userID string, keyword string) ([]*entity.Subject, error)
}
