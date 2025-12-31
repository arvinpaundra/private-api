package repository

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/submission/entity"
)

type SubmissionWriter interface {
	Save(ctx context.Context, submission *entity.Submission) error
}
