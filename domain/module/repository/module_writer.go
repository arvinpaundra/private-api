package repository

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/module/entity"
)

type ModuleWriter interface {
	Save(ctx context.Context, module *entity.Module) error
}
