package service

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/shared/interfaces"
	"github.com/arvinpaundra/private-api/domain/subject/repository"
)

type GetSubjectNamesCommand struct {
	SubjectIDs []string `validate:"required"`
}

type GetSubjectNames struct {
	authStorage   interfaces.AuthenticatedUser
	subjectReader repository.SubjectReader
}

func NewGetSubjectNames(
	authStorage interfaces.AuthenticatedUser,
	subjectReader repository.SubjectReader,
) *GetSubjectNames {
	return &GetSubjectNames{
		authStorage:   authStorage,
		subjectReader: subjectReader,
	}
}

func (s *GetSubjectNames) Execute(ctx context.Context, command *GetSubjectNamesCommand) (map[string]string, error) {
	names := make(map[string]string)

	for _, subjectID := range command.SubjectIDs {
		subject, err := s.subjectReader.FindSubjectByID(ctx, subjectID, s.authStorage.GetUserId())
		if err != nil {
			// Skip if not found, continue with others
			continue
		}
		names[subjectID] = subject.Name
	}

	return names, nil
}
