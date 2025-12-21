package repository

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/subject/entity"
)

type SubjectWriter interface {
	Save(ctx context.Context, subject *entity.Subject) error
}
