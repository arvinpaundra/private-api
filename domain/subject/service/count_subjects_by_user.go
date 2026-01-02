package service

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/subject/repository"
)

type CountSubjectsByUser struct {
	subjectReader repository.SubjectReader
}

func NewCountSubjectsByUser(subjectReader repository.SubjectReader) *CountSubjectsByUser {
	return &CountSubjectsByUser{
		subjectReader: subjectReader,
	}
}

func (s *CountSubjectsByUser) Execute(ctx context.Context, userID string) (int, error) {
	count, err := s.subjectReader.CountByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}

	return count, nil
}
