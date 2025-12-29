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

func (a *SubjectACLAdapter) GetSubjectName(ctx context.Context, subjectID string, userID string) (string, error) {
	subjectService := service.NewFindDetailSubject(
		a.authStorage,
		subject.NewSubjectReaderRepository(a.db),
	)

	subjectDetail, err := subjectService.Execute(ctx, &service.FindDetailSubjectCommand{
		ID: subjectID,
	})
	if err != nil {
		return "", err
	}

	return subjectDetail.Name, nil
}

func (a *SubjectACLAdapter) GetSubjectNames(ctx context.Context, subjectIDs []string, userID string) (map[string]string, error) {
	subjectService := service.NewGetSubjectNames(
		a.authStorage,
		subject.NewSubjectReaderRepository(a.db),
	)

	names, err := subjectService.Execute(ctx, &service.GetSubjectNamesCommand{
		SubjectIDs: subjectIDs,
	})
	if err != nil {
		return nil, err
	}

	return names, nil
}
