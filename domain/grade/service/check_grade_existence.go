package service

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/grade/repository"
	"github.com/arvinpaundra/private-api/domain/shared/interfaces"
)

type CheckGradeExistenceCommand struct {
	GradeID string `validate:"required"`
}

type CheckGradeExistence struct {
	authStorage interfaces.AuthenticatedUser
	gradeReader repository.GradeReader
}

func NewCheckGradeExistence(
	authStorage interfaces.AuthenticatedUser,
	gradeReader repository.GradeReader,
) *CheckGradeExistence {
	return &CheckGradeExistence{
		authStorage: authStorage,
		gradeReader: gradeReader,
	}
}

func (s *CheckGradeExistence) Execute(ctx context.Context, command *CheckGradeExistenceCommand) (bool, error) {
	exists, err := s.gradeReader.IsGradeExist(ctx, command.GradeID, s.authStorage.GetUserId())
	if err != nil {
		return false, err
	}

	if exists {
		return true, nil
	}

	return false, nil
}
