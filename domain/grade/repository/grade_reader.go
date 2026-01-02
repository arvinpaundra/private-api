package repository

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/grade/entity"
)

type GradeReader interface {
	HasSimilarGrade(ctx context.Context, name string, userID string) (bool, error)
	HasSimilarGradeExclusive(ctx context.Context, name string, userID string, excludeGradeID string) (bool, error)
	IsGradeExist(ctx context.Context, gradeID string, userID string) (bool, error)
	FindGradeByID(ctx context.Context, gradeID string, userID string) (*entity.Grade, error)
	CountByUserID(ctx context.Context, userID string) (int, error)
	AllGrades(ctx context.Context, userID string, keyword string) ([]*entity.Grade, error)
}
