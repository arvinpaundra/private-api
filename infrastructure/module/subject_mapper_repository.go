package module

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/module/repository"
	"github.com/arvinpaundra/private-api/domain/shared/interfaces"
	"github.com/arvinpaundra/private-api/domain/subject/service"
	"github.com/arvinpaundra/private-api/infrastructure/subject"
	"gorm.io/gorm"
)

var _ repository.SubjectACL = (*SubjectACLAdapter)(nil)

type SubjectACLAdapter struct {
	db          *gorm.DB
	authStorage interfaces.AuthenticatedUser
}

func NewSubjectACLAdapter(db *gorm.DB, authStorage interfaces.AuthenticatedUser) *SubjectACLAdapter {
	return &SubjectACLAdapter{
		db:          db,
		authStorage: authStorage,
	}
}

func (a *SubjectACLAdapter) IsSubjectExist(ctx context.Context, subjectID string, userID string) (bool, error) {
	subjectService := service.NewCheckSubjectExistence(
		a.authStorage,
		subject.NewSubjectReaderRepository(a.db),
	)

	exists, err := subjectService.Execute(ctx, &service.CheckSubjectExistenceCommand{
		SubjectID: subjectID,
	})
	if err != nil {
		return false, err
	}

	return exists, nil
}
