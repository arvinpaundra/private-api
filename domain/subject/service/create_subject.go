package service

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/shared/interfaces"
	"github.com/arvinpaundra/private-api/domain/subject/constant"
	"github.com/arvinpaundra/private-api/domain/subject/entity"
	"github.com/arvinpaundra/private-api/domain/subject/repository"
)

type CreateSubjectCommand struct {
	Name        string  `json:"name" validate:"required,max=100"`
	Description *string `json:"description,omitempty"`
}

type CreateSubject struct {
	authStorage   interfaces.AuthenticatedUser
	subjectReader repository.SubjectReader
	subjectWriter repository.SubjectWriter
}

func NewCreateSubject(
	authStorage interfaces.AuthenticatedUser,
	subjectReader repository.SubjectReader,
	subjectWriter repository.SubjectWriter,
) *CreateSubject {
	return &CreateSubject{
		authStorage:   authStorage,
		subjectReader: subjectReader,
		subjectWriter: subjectWriter,
	}
}

func (s *CreateSubject) Execute(ctx context.Context, command *CreateSubjectCommand) error {
	// Check if there have similar subject
	hasSimilarSubject, err := s.subjectReader.HasSimilarSubject(ctx, command.Name, s.authStorage.GetUserId())
	if err != nil {
		return err
	}

	if hasSimilarSubject {
		return constant.ErrSubjectAlreadyExists
	}

	// Create subject
	subject := entity.NewSubject(s.authStorage.GetUserId(), command.Name, command.Description)

	// Store subject to persistent storage
	err = s.subjectWriter.Save(ctx, subject)
	if err != nil {
		return err
	}

	return nil
}
