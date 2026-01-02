package dashboard

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/dashboard/repository"
	"github.com/arvinpaundra/private-api/model"
	"gorm.io/gorm"
)

var _ repository.SubjectACL = (*SubjectACLAdapter)(nil)

type SubjectACLAdapter struct {
	db *gorm.DB
}

func NewSubjectACLAdapter(db *gorm.DB) *SubjectACLAdapter {
	return &SubjectACLAdapter{
		db: db,
	}
}

func (a *SubjectACLAdapter) CountSubjectsByUserID(ctx context.Context, userID string) (int, error) {
	var count int64

	err := a.db.Model(&model.Subject{}).
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
