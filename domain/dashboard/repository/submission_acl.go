package repository

import "context"

type SubmissionACL interface {
	CountSubmittedSubmissions(ctx context.Context) (int, error)
}
