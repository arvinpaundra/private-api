package service

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/module/repository"
	"github.com/arvinpaundra/private-api/domain/module/response"
	"github.com/arvinpaundra/private-api/domain/shared/interfaces"
)

type FindDetailModuleCommand struct {
	Slug string `form:"slug"`
}

type FindDetailModule struct {
	authStorage  interfaces.AuthenticatedUser
	moduleReader repository.ModuleReader
	subjectACL   repository.SubjectACL
	gradeACL     repository.GradeACL
}

func NewFindDetailModule(
	authStorage interfaces.AuthenticatedUser,
	moduleReader repository.ModuleReader,
	subjectACL repository.SubjectACL,
	gradeACL repository.GradeACL,
) *FindDetailModule {
	return &FindDetailModule{
		authStorage:  authStorage,
		moduleReader: moduleReader,
		subjectACL:   subjectACL,
		gradeACL:     gradeACL,
	}
}

func (s *FindDetailModule) Execute(ctx context.Context, command *FindDetailModuleCommand) (*response.ModuleDetail, error) {
	module, err := s.moduleReader.FindModuleDetailBySlug(ctx, command.Slug, s.authStorage.GetUserId())
	if err != nil {
		return nil, err
	}

	// Get subject name via ACL
	subjectName, err := s.subjectACL.GetSubjectName(ctx, module.SubjectID, s.authStorage.GetUserId())
	if err != nil {
		return nil, err
	}

	// Get grade name via ACL
	gradeName, err := s.gradeACL.GetGradeName(ctx, module.GradeID, s.authStorage.GetUserId())
	if err != nil {
		return nil, err
	}

	// Transform questions and choices
	questions := make([]*response.Question, len(module.Questions))

	for i, question := range module.Questions {
		choices := make([]*response.ChoiceWithAnswer, len(question.Choices))

		for j, choice := range question.Choices {
			choices[j] = &response.ChoiceWithAnswer{
				ID:              choice.ID,
				Content:         choice.Content,
				IsCorrectAnswer: choice.IsCorrectAnswer,
			}
		}

		questions[i] = &response.Question{
			ID:      question.ID,
			Content: question.Content,
			Slug:    question.Slug,
			Choices: choices,
		}
	}

	return &response.ModuleDetail{
		ID:          module.ID,
		Title:       module.Title,
		Slug:        module.Slug,
		Description: module.Description,
		Type:        module.Type,
		IsPublished: module.IsPublished,
		Subject: &response.Subject{
			ID:   module.SubjectID,
			Name: subjectName,
		},
		Grade: &response.Grade{
			ID:   module.GradeID,
			Name: gradeName,
		},
		Questions: questions,
	}, nil
}
