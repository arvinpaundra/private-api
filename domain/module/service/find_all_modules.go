package service

import (
	"context"
	"slices"

	"github.com/arvinpaundra/private-api/core/format"
	"github.com/arvinpaundra/private-api/domain/module/repository"
	"github.com/arvinpaundra/private-api/domain/module/response"
	"github.com/arvinpaundra/private-api/domain/shared/interfaces"
)

type FindAllModulesCommand struct {
	Keyword   string `form:"keyword"`
	SubjectID string `form:"subject_id"`
	GradeID   string `form:"grade_id"`
	Page      int    `form:"page"`
	PerPage   int    `form:"per_page"`
}

type FindAllModulesResponse struct {
	Modules    []*response.Module `json:"modules"`
	Pagination format.Pagination  `json:"pagination"`
}

type FindAllModules struct {
	authStorage  interfaces.AuthenticatedUser
	moduleReader repository.ModuleReader
	subjectACL   repository.SubjectACL
	gradeACL     repository.GradeACL
}

func NewFindAllModules(
	authStorage interfaces.AuthenticatedUser,
	moduleReader repository.ModuleReader,
	subjectACL repository.SubjectACL,
	gradeACL repository.GradeACL,
) *FindAllModules {
	return &FindAllModules{
		authStorage:  authStorage,
		moduleReader: moduleReader,
		subjectACL:   subjectACL,
		gradeACL:     gradeACL,
	}
}

func (s *FindAllModules) Execute(ctx context.Context, command *FindAllModulesCommand) (*FindAllModulesResponse, error) {
	// Validate pagination values
	page := format.ValidatePage(command.Page)
	perPage := format.ValidatePerPage(command.PerPage)
	offset := format.CalculateOffset(command.Page, command.PerPage)

	// Get total count
	total, err := s.moduleReader.TotalModules(ctx,
		s.authStorage.GetUserId(),
		command.SubjectID,
		command.GradeID,
		command.Keyword,
	)
	if err != nil {
		return nil, err
	}

	// Get modules
	modules, err := s.moduleReader.FindAllModules(ctx,
		s.authStorage.GetUserId(),
		command.SubjectID,
		command.GradeID,
		command.Keyword,
		perPage,
		offset,
	)
	if err != nil {
		return nil, err
	}

	// Extract unique subject and grade IDs from modules
	var subjectIDs, gradeIDs []string

	for _, module := range modules {
		if !slices.Contains(subjectIDs, module.SubjectID) {
			subjectIDs = append(subjectIDs, module.SubjectID)
		}

		if !slices.Contains(gradeIDs, module.GradeID) {
			gradeIDs = append(gradeIDs, module.GradeID)
		}
	}

	// Get all subject names
	subjectNames, err := s.subjectACL.GetSubjectNames(ctx, subjectIDs, s.authStorage.GetUserId())
	if err != nil {
		return nil, err
	}

	// Get all grade names
	gradeNames, err := s.gradeACL.GetGradeNames(ctx, gradeIDs, s.authStorage.GetUserId())
	if err != nil {
		return nil, err
	}

	// Transform to response with relations
	results := make([]*response.Module, len(modules))

	for i, module := range modules {
		results[i] = &response.Module{
			ID:             module.ID,
			UserID:         module.UserID,
			SubjectID:      module.SubjectID,
			GradeID:        module.GradeID,
			Title:          module.Title,
			Slug:           module.Slug,
			Description:    module.Description,
			Type:           module.Type,
			IsPublished:    module.IsPublished,
			QuestionsCount: len(module.Questions),
			Subject: &response.Subject{
				ID:   module.SubjectID,
				Name: subjectNames[module.SubjectID],
			},
			Grade: &response.Grade{
				ID:   module.GradeID,
				Name: gradeNames[module.GradeID],
			},
		}
	}

	// Create pagination
	pagination := format.NewPagination(page, perPage, total)

	return &FindAllModulesResponse{
		Modules:    results,
		Pagination: pagination,
	}, nil
}
