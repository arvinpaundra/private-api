package service

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/shared/interfaces"
	"github.com/arvinpaundra/private-api/domain/subject/repository"
	"github.com/arvinpaundra/private-api/domain/subject/response"
)

type FindAllSubjectsCommand struct {
	Keyword string `form:"keyword"`
}

type FindAllSubjects struct {
	authStorage   interfaces.AuthenticatedUser
	subjectReader repository.SubjectReader
}

func NewFindAllSubjects(
	authStorage interfaces.AuthenticatedUser,
	subjectReader repository.SubjectReader,
) *FindAllSubjects {
	return &FindAllSubjects{
		authStorage:   authStorage,
		subjectReader: subjectReader,
	}
}

func (s *FindAllSubjects) Execute(ctx context.Context, command *FindAllSubjectsCommand) ([]*response.Subject, error) {
	subjects, err := s.subjectReader.AllSubjects(ctx, s.authStorage.GetUserId(), command.Keyword)
	if err != nil {
		return nil, err
	}

	results := make([]*response.Subject, len(subjects))

	for i, subject := range subjects {
		results[i] = &response.Subject{
			ID:          subject.ID,
			Name:        subject.Name,
			Description: subject.Description,
		}
	}

	return results, nil
}
