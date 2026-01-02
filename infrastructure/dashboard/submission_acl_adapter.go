package dashboard

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/dashboard/repository"
	"github.com/arvinpaundra/private-api/domain/submission/service"
	"github.com/arvinpaundra/private-api/infrastructure/submission"
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
	svc := service.NewCountSubmittedSubmissions(
		submission.NewSubmissionReaderRepository(a.db),
	)

	count, err := svc.Execute(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}
