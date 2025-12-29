package service

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/shared/interfaces"
	"github.com/arvinpaundra/private-api/domain/subject/repository"
)

type CheckSubjectExistenceCommand struct {
	SubjectID string `validate:"required"`
}

type CheckSubjectExistence struct {
	authStorage   interfaces.AuthenticatedUser
	subjectReader repository.SubjectReader
}

func NewCheckSubjectExistence(
	authStorage interfaces.AuthenticatedUser,
	subjectReader repository.SubjectReader,
) *CheckSubjectExistence {
	return &CheckSubjectExistence{
		authStorage:   authStorage,
		subjectReader: subjectReader,
	}
}

func (s *CheckSubjectExistence) Execute(ctx context.Context, command *CheckSubjectExistenceCommand) (bool, error) {
	exists, err := s.subjectReader.IsSubjectExist(ctx, command.SubjectID, s.authStorage.GetUserId())
	if err != nil {
		return false, err
	}

	if exists {
		return true, nil
	}

	return false, nil
}
