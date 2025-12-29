package module

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/grade/service"
	"github.com/arvinpaundra/private-api/domain/module/repository"
	"github.com/arvinpaundra/private-api/domain/shared/interfaces"
	"github.com/arvinpaundra/private-api/infrastructure/grade"
	"gorm.io/gorm"
)

var _ repository.GradeACL = (*GradeACLAdapter)(nil)

type GradeACLAdapter struct {
	db          *gorm.DB
	authStorage interfaces.AuthenticatedUser
}

func NewGradeACLAdapter(db *gorm.DB, authStorage interfaces.AuthenticatedUser) *GradeACLAdapter {
	return &GradeACLAdapter{
		db:          db,
		authStorage: authStorage,
	}
}

func (a *GradeACLAdapter) IsGradeExist(ctx context.Context, gradeID string, userID string) (bool, error) {
	gradeService := service.NewCheckGradeExistence(
		a.authStorage,
		grade.NewGradeReaderRepository(a.db),
	)

	exists, err := gradeService.Execute(ctx, &service.CheckGradeExistenceCommand{
		GradeID: gradeID,
	})
	if err != nil {
		return false, err
	}

	return exists, nil
}
