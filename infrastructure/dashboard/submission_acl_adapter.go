package dashboard

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/dashboard/repository"
	"github.com/arvinpaundra/private-api/model"
	"gorm.io/gorm"
)

var _ repository.SubmissionACL = (*SubmissionACLAdapter)(nil)

type SubmissionACLAdapter struct {
	db *gorm.DB
}

func NewSubmissionACLAdapter(db *gorm.DB) *SubmissionACLAdapter {
	return &SubmissionACLAdapter{
		db: db,
	}
}

func (a *SubmissionACLAdapter) CountSubmittedSubmissions(ctx context.Context) (int, error) {
	var count int64

	err := a.db.Model(&model.Submission{}).
		WithContext(ctx).
		Where("status = ?", model.Submitted).
		Count(&count).
		Error

	if err != nil {
		return 0, err
	}

	return int(count), nil
}
