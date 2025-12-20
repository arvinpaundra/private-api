package repository

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/auth/entity"
)

type UserWriter interface {
	Save(ctx context.Context, user *entity.User) error
}
