package repository

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/user/entity"
)

type UserReader interface {
	FindByID(ctx context.Context, userID string) (*entity.User, error)
}
