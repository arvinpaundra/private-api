package service

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/grade/repository"
	"github.com/arvinpaundra/private-api/domain/shared/interfaces"
)

type DeleteGradeCommand struct {
	ID string `json:"id" validate:"required"`
}

type DeleteGrade struct {
	authStorage interfaces.AuthenticatedUser
	gradeReader repository.GradeReader
	gradeWriter repository.GradeWriter
}

func NewDeleteGrade(
	authStorage interfaces.AuthenticatedUser,
	gradeReader repository.GradeReader,
	gradeWriter repository.GradeWriter,
) *DeleteGrade {
	return &DeleteGrade{
		authStorage: authStorage,
		gradeReader: gradeReader,
		gradeWriter: gradeWriter,
	}
}

func (s *DeleteGrade) Execute(ctx context.Context, command *DeleteGradeCommand) error {
	// Check if grade exists
	grade, err := s.gradeReader.FindGradeByID(ctx, command.ID, s.authStorage.GetUserId())
	if err != nil {
		return err
	}

	// Delete grade
	grade.MarkRemove()

	// Update grade as deleted in persistent storage
	err = s.gradeWriter.Save(ctx, grade)
	if err != nil {
		return err
	}

	return nil
}
