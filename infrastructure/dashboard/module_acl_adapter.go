package dashboard

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/dashboard/repository"
	"github.com/arvinpaundra/private-api/model"
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
	var count int64

	err := a.db.Model(&model.Module{}).
		WithContext(ctx).
		Where("user_id = ?", userID).
		Where("deleted_at IS NULL").
		Count(&count).
		Error

	if err != nil {
		return 0, err
	}

	return int(count), nil
}
