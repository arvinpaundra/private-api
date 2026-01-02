package dashboard

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/dashboard/repository"
	"github.com/arvinpaundra/private-api/domain/subject/service"
	"github.com/arvinpaundra/private-api/infrastructure/subject"
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
	svc := service.NewCountSubjectsByUser(
		subject.NewSubjectReaderRepository(a.db),
	)

	count, err := svc.Execute(ctx, userID)
	if err != nil {
		return 0, err
	}

	return count, nil
}
