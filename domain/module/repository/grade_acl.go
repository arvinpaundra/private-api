package repository

import (
	"context"
)

type GradeACL interface {
	IsGradeExist(ctx context.Context, gradeID string, userID string) (bool, error)
}
