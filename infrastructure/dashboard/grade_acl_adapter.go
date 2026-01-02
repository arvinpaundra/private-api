package dashboard

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/dashboard/repository"
	"github.com/arvinpaundra/private-api/domain/grade/service"
	"github.com/arvinpaundra/private-api/infrastructure/grade"
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
	svc := service.NewCountGradesByUser(
		grade.NewGradeReaderRepository(a.db),
	)

	count, err := svc.Execute(ctx, userID)
	if err != nil {
		return 0, err
	}

	return count, nil
}
