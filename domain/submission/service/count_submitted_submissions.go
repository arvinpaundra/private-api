package service

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/submission/repository"
)

type CountSubmittedSubmissions struct {
	submissionReader repository.SubmissionReader
}

func NewCountSubmittedSubmissions(submissionReader repository.SubmissionReader) *CountSubmittedSubmissions {
	return &CountSubmittedSubmissions{
		submissionReader: submissionReader,
	}
}

func (s *CountSubmittedSubmissions) Execute(ctx context.Context) (int, error) {
	count, err := s.submissionReader.CountSubmitted(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}
