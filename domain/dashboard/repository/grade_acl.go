package repository

import "context"

type GradeACL interface {
	CountGradesByUserID(ctx context.Context, userID string) (int, error)
}
