package service

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/grade/repository"
)

type CountGradesByUser struct {
	gradeReader repository.GradeReader
}

func NewCountGradesByUser(gradeReader repository.GradeReader) *CountGradesByUser {
	return &CountGradesByUser{
		gradeReader: gradeReader,
	}
}

func (s *CountGradesByUser) Execute(ctx context.Context, userID string) (int, error) {
	count, err := s.gradeReader.CountByUserID(ctx, userID)
	if err != nil {
		return 0, err
	}

	return count, nil
}
