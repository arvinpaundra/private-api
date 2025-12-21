package service

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/grade/repository"
	"github.com/arvinpaundra/private-api/domain/grade/response"
	"github.com/arvinpaundra/private-api/domain/shared/interfaces"
)

type FindAllGradesCommand struct {
	Keyword string `form:"keyword"`
}

type FindAllGrades struct {
	authStorage interfaces.AuthenticatedUser
	gradeReader repository.GradeReader
}

func NewFindAllGrades(
	authStorage interfaces.AuthenticatedUser,
	gradeReader repository.GradeReader,
) *FindAllGrades {
	return &FindAllGrades{
		authStorage: authStorage,
		gradeReader: gradeReader,
	}
}

func (s *FindAllGrades) Execute(ctx context.Context, command *FindAllGradesCommand) ([]*response.Grade, error) {
	grades, err := s.gradeReader.AllGrades(ctx, s.authStorage.GetUserId(), command.Keyword)
	if err != nil {
		return nil, err
	}

	results := make([]*response.Grade, len(grades))

	for i, grade := range grades {
		results[i] = &response.Grade{
			ID:          grade.ID,
			Name:        grade.Name,
			Description: grade.Description,
		}
	}

	return results, nil
}
