package dashboard

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/dashboard/repository"
	"github.com/arvinpaundra/private-api/model"
	"gorm.io/gorm"
)

var _ repository.GradeACL = (*GradeACLAdapter)(nil)

type GradeACLAdapter struct {
	db *gorm.DB
}

func NewGradeACLAdapter(db *gorm.DB) *GradeACLAdapter {
	return &GradeACLAdapter{
		db: db,
	}
}

func (a *GradeACLAdapter) CountGradesByUserID(ctx context.Context, userID string) (int, error) {
	var count int64

	err := a.db.Model(&model.Grade{}).
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
