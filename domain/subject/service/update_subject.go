package service

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/shared/interfaces"
	"github.com/arvinpaundra/private-api/domain/subject/constant"
	"github.com/arvinpaundra/private-api/domain/subject/repository"
)

type UpdateSubjectCommand struct {
	ID          string  `form:"id"`
	Name        string  `json:"name" validate:"required,max=100"`
	Description *string `json:"description"`
}

type UpdateSubject struct {
	authStorage   interfaces.AuthenticatedUser
	subjectReader repository.SubjectReader
	subjectWriter repository.SubjectWriter
}

func NewUpdateSubject(
	authStorage interfaces.AuthenticatedUser,
	subjectReader repository.SubjectReader,
	subjectWriter repository.SubjectWriter,
) *UpdateSubject {
	return &UpdateSubject{
		authStorage:   authStorage,
		subjectReader: subjectReader,
		subjectWriter: subjectWriter,
	}
}

func (s *UpdateSubject) Execute(ctx context.Context, command *UpdateSubjectCommand) error {
	// Check if subject exists
	subject, err := s.subjectReader.FindSubjectByID(ctx, command.ID, s.authStorage.GetUserId())
	if err != nil {
		return err
	}

	// Check if there have similar subject excluding current subject
	hasSimilarSubject, err := s.subjectReader.HasSimilarSubjectExclusive(ctx, command.Name, s.authStorage.GetUserId(), command.ID)
	if err != nil {
		return err
	}

	if hasSimilarSubject {
		return constant.ErrSubjectAlreadyExists
	}

	// Update subject
	subject.Update(command.Name, command.Description)

	// Store updated subject to persistent storage
	err = s.subjectWriter.Save(ctx, subject)
	if err != nil {
		return err
	}

	return nil
}
