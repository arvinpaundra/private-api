package service

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/grade/repository"
	"github.com/arvinpaundra/private-api/domain/shared/interfaces"
)

type GetGradeNamesCommand struct {
	GradeIDs []string `validate:"required"`
}

type GetGradeNames struct {
	authStorage interfaces.AuthenticatedUser
	gradeReader repository.GradeReader
}

func NewGetGradeNames(
	authStorage interfaces.AuthenticatedUser,
	gradeReader repository.GradeReader,
) *GetGradeNames {
	return &GetGradeNames{
		authStorage: authStorage,
		gradeReader: gradeReader,
	}
}

func (s *GetGradeNames) Execute(ctx context.Context, command *GetGradeNamesCommand) (map[string]string, error) {
	names := make(map[string]string)

	for _, gradeID := range command.GradeIDs {
		grade, err := s.gradeReader.FindGradeByID(ctx, gradeID, s.authStorage.GetUserId())
		if err != nil {
			// Skip if not found, continue with others
			continue
		}
		names[gradeID] = grade.Name
	}

	return names, nil
}
