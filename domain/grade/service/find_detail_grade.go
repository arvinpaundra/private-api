package service

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/grade/repository"
	"github.com/arvinpaundra/private-api/domain/grade/response"
	"github.com/arvinpaundra/private-api/domain/shared/interfaces"
)

type FindDetailGradeCommand struct {
	ID string `form:"id"`
}

type FindDetailGrade struct {
	authStorage interfaces.AuthenticatedUser
	gradeReader repository.GradeReader
}

func NewFindDetailGrade(
	authStorage interfaces.AuthenticatedUser,
	gradeReader repository.GradeReader,
) *FindDetailGrade {
	return &FindDetailGrade{
		authStorage: authStorage,
		gradeReader: gradeReader,
	}
}

func (s *FindDetailGrade) Execute(ctx context.Context, command *FindDetailGradeCommand) (*response.Grade, error) {
	grade, err := s.gradeReader.FindGradeByID(ctx, command.ID, s.authStorage.GetUserId())
	if err != nil {
		return nil, err
	}

	return &response.Grade{
		ID:          grade.ID,
		Name:        grade.Name,
		Description: grade.Description,
	}, nil
}
