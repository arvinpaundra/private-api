package service

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/shared/interfaces"
	"github.com/arvinpaundra/private-api/domain/subject/repository"
)

type DeleteSubjectCommand struct {
	ID string `json:"id" validate:"required"`
}

type DeleteSubject struct {
	authStorage   interfaces.AuthenticatedUser
	subjectReader repository.SubjectReader
	subjectWriter repository.SubjectWriter
}

func NewDeleteSubject(
	authStorage interfaces.AuthenticatedUser,
	subjectReader repository.SubjectReader,
	subjectWriter repository.SubjectWriter,
) *DeleteSubject {
	return &DeleteSubject{
		authStorage:   authStorage,
		subjectReader: subjectReader,
		subjectWriter: subjectWriter,
	}
}

func (s *DeleteSubject) Execute(ctx context.Context, command *DeleteSubjectCommand) error {
	// Check if subject exists
	subject, err := s.subjectReader.FindSubjectByID(ctx, command.ID, s.authStorage.GetUserId())
	if err != nil {
		return err
	}

	// Delete subject
	subject.MarkRemove()

	// Update subject as deleted in persistent storage
	err = s.subjectWriter.Save(ctx, subject)
	if err != nil {
		return err
	}

	return nil
}
