package service

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/grade/constant"
	"github.com/arvinpaundra/private-api/domain/grade/entity"
	"github.com/arvinpaundra/private-api/domain/grade/repository"
	"github.com/arvinpaundra/private-api/domain/shared/interfaces"
)

type CreateGradeCommand struct {
	Name        string  `json:"name" validate:"required,max=100"`
	Description *string `json:"description,omitempty"`
}

type CreateGrade struct {
	authStorage interfaces.AuthenticatedUser
	gradeReader repository.GradeReader
	gradeWriter repository.GradeWriter
}

func NewCreateGrade(
	authStorage interfaces.AuthenticatedUser,
	gradeReader repository.GradeReader,
	gradeWriter repository.GradeWriter,
) *CreateGrade {
	return &CreateGrade{
		authStorage: authStorage,
		gradeReader: gradeReader,
		gradeWriter: gradeWriter,
	}
}

func (s *CreateGrade) Execute(ctx context.Context, command *CreateGradeCommand) error {
	// Check if there have similar grade
	hasSimilarGrade, err := s.gradeReader.HasSimilarGrade(ctx, command.Name, s.authStorage.GetUserId())
	if err != nil {
		return err
	}

	if hasSimilarGrade {
		return constant.ErrGradeAlreadyExists
	}

	// Create grade
	grade := entity.NewGrade(s.authStorage.GetUserId(), command.Name, command.Description)

	// Store grade to persistent storage
	err = s.gradeWriter.Save(ctx, grade)
	if err != nil {
		return err
	}

	return nil
}
