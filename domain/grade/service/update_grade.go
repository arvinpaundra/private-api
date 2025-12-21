package service

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/grade/constant"
	"github.com/arvinpaundra/private-api/domain/grade/repository"
	"github.com/arvinpaundra/private-api/domain/shared/interfaces"
)

type UpdateGradeCommand struct {
	ID          string  `form:"id"`
	Name        string  `json:"name" validate:"required,max=100"`
	Description *string `json:"description"`
}

type UpdateGrade struct {
	authStorage interfaces.AuthenticatedUser
	gradeReader repository.GradeReader
	gradeWriter repository.GradeWriter
}

func NewUpdateGrade(
	authStorage interfaces.AuthenticatedUser,
	gradeReader repository.GradeReader,
	gradeWriter repository.GradeWriter,
) *UpdateGrade {
	return &UpdateGrade{
		authStorage: authStorage,
		gradeReader: gradeReader,
		gradeWriter: gradeWriter,
	}
}

func (s *UpdateGrade) Execute(ctx context.Context, command *UpdateGradeCommand) error {
	// Check if grade exists
	grade, err := s.gradeReader.FindGradeByID(ctx, command.ID, s.authStorage.GetUserId())
	if err != nil {
		return err
	}

	// Check if there have similar grade excluding current grade
	hasSimilarGrade, err := s.gradeReader.HasSimilarGradeExclusive(ctx, command.Name, s.authStorage.GetUserId(), command.ID)
	if err != nil {
		return err
	}

	if hasSimilarGrade {
		return constant.ErrGradeAlreadyExists
	}

	// Update grade
	grade.Update(command.Name, command.Description)

	// Store updated grade to persistent storage
	err = s.gradeWriter.Save(ctx, grade)
	if err != nil {
		return err
	}

	return nil
}
