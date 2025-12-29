package repository

import (
	"context"
)

type SubjectACL interface {
	IsSubjectExist(ctx context.Context, subjectID string, userID string) (bool, error)
	GetSubjectName(ctx context.Context, subjectID string, userID string) (string, error)
	GetSubjectNames(ctx context.Context, subjectIDs []string, userID string) (map[string]string, error)
}
