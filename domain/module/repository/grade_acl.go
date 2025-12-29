package repository

import (
	"context"
)

type GradeACL interface {
	IsGradeExist(ctx context.Context, gradeID string, userID string) (bool, error)
	GetGradeName(ctx context.Context, gradeID string, userID string) (string, error)
	GetGradeNames(ctx context.Context, gradeIDs []string, userID string) (map[string]string, error)
}
