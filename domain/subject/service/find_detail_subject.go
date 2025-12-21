package service

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/shared/interfaces"
	"github.com/arvinpaundra/private-api/domain/subject/repository"
	"github.com/arvinpaundra/private-api/domain/subject/response"
)

type FindDetailSubjectCommand struct {
	ID string `form:"id"`
}

type FindDetailSubject struct {
	authStorage   interfaces.AuthenticatedUser
	subjectReader repository.SubjectReader
}

func NewFindDetailSubject(
	authStorage interfaces.AuthenticatedUser,
	subjectReader repository.SubjectReader,
) *FindDetailSubject {
	return &FindDetailSubject{
		authStorage:   authStorage,
		subjectReader: subjectReader,
	}
}

func (s *FindDetailSubject) Execute(ctx context.Context, command *FindDetailSubjectCommand) (*response.Subject, error) {
	subject, err := s.subjectReader.FindSubjectByID(ctx, command.ID, s.authStorage.GetUserId())
	if err != nil {
		return nil, err
	}

	return &response.Subject{
		ID:          subject.ID,
		Name:        subject.Name,
		Description: subject.Description,
	}, nil
}
