package service

import (
	"context"
	"time"

	"github.com/arvinpaundra/private-api/domain/submission/repository"
	"github.com/arvinpaundra/private-api/domain/submission/response"
)

type FindAllSubmissionQuery struct {
	Keyword string `form:"keyword"`
}

type FindAllSubmission struct {
	submissionReader repository.SubmissionReader
	moduleACL        repository.ModuleACL
}

func NewFindAllSubmission(
	submissionReader repository.SubmissionReader,
	moduleACL repository.ModuleACL,
) *FindAllSubmission {
	return &FindAllSubmission{
		submissionReader: submissionReader,
		moduleACL:        moduleACL,
	}
}

func (s *FindAllSubmission) Execute(ctx context.Context, query *FindAllSubmissionQuery) ([]*response.ModuleSubmissionGroup, error) {
	// Get all published modules
	modules, err := s.moduleACL.GetAllPublishedModules(ctx, query.Keyword)
	if err != nil {
		return nil, err
	}

	// Early return if no modules found
	if len(modules) == 0 {
		return []*response.ModuleSubmissionGroup{}, nil
	}

	// Extract module IDs for bulk query
	moduleIDs := make([]string, len(modules))
	for i, module := range modules {
		moduleIDs[i] = module.ID
	}

	// Fetch all submissions grouped by module in a single query
	submissionsByModule, err := s.submissionReader.FindAllSubmittedGroupedByModule(ctx, moduleIDs)
	if err != nil {
		return nil, err
	}

	// Build the grouped result
	result := make([]*response.ModuleSubmissionGroup, 0, len(modules))

	for _, module := range modules {
		submissions := submissionsByModule[module.ID]

		// Map submissions to summary
		submissionSummaries := make([]*response.SubmissionSummary, len(submissions))
		for i, submission := range submissions {
			submittedAt := ""

			if submission.SubmittedAt != nil {
				submittedAt = submission.SubmittedAt.Format(time.DateTime)
			}

			submissionSummaries[i] = &response.SubmissionSummary{
				StudentName:    submission.StudentName,
				TotalCorrect:   submission.Score(),
				TotalQuestions: submission.TotalQuestions,
				SubmittedAt:    submittedAt,
			}
		}

		// Build module submission group
		group := &response.ModuleSubmissionGroup{
			Module: &response.ModuleWithRelations{
				ID:    module.ID,
				Title: module.Title,
				Slug:  module.Slug,
				Grade: &response.Grade{
					ID:   module.Grade.ID,
					Name: module.Grade.Name,
				},
				Subject: &response.Subject{
					ID:   module.Subject.ID,
					Name: module.Subject.Name,
				},
			},
			TotalSubmissions: len(submissions),
			Submissions:      submissionSummaries,
		}

		result = append(result, group)
	}

	return result, nil
}
