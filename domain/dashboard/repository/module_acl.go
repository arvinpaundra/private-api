package repository

import "context"

type ModuleACL interface {
	CountModulesByUserID(ctx context.Context, userID string) (int, error)
}
