package repository

import "context"

type SubjectACL interface {
	CountSubjectsByUserID(ctx context.Context, userID string) (int, error)
}
