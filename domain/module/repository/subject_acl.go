package repository

import (
	"context"
)

type SubjectACL interface {
	IsSubjectExist(ctx context.Context, subjectID string, userID string) (bool, error)
}
