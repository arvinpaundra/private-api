package dashboard

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/dashboard/repository"
	"github.com/arvinpaundra/private-api/domain/module/service"
	"github.com/arvinpaundra/private-api/infrastructure/module"
	"gorm.io/gorm"
)

var _ repository.ModuleACL = (*ModuleACLAdapter)(nil)

type ModuleACLAdapter struct {
	db *gorm.DB
}

func NewModuleACLAdapter(db *gorm.DB) *ModuleACLAdapter {
	return &ModuleACLAdapter{
		db: db,
	}
}

func (a *ModuleACLAdapter) CountModulesByUserID(ctx context.Context, userID string) (int, error) {
	svc := service.NewCountModulesByUser(
		module.NewModuleReaderRepository(a.db),
	)

	count, err := svc.Execute(ctx, userID)
	if err != nil {
		return 0, err
	}

	return count, nil
}
