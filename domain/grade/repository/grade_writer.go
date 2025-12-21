package repository

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/grade/entity"
)

type GradeWriter interface {
	Save(ctx context.Context, grade *entity.Grade) error
}
