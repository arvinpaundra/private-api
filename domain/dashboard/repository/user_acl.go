package repository

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/dashboard/entity"
)

type UserACL interface {
	GetUserByID(ctx context.Context, userID string) (*entity.User, error)
}
