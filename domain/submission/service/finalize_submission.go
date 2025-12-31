package service

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/submission/constant"
	"github.com/arvinpaundra/private-api/domain/submission/repository"
	"github.com/arvinpaundra/private-api/domain/submission/response"
)

type FinalizeSubmissionCommand struct {
	SubmissionCode string `json:"-" validate:"required"`
	ModuleSlug     string `json:"-" validate:"required"`
}

type FinalizeSubmission struct {
	submissionReader repository.SubmissionReader
	moduleACL        repository.ModuleACL
	uow              repository.UnitOfWork
}

func NewFinalizeSubmission(
	submissionReader repository.SubmissionReader,
	moduleACL repository.ModuleACL,
	uow repository.UnitOfWork,
) *FinalizeSubmission {
	return &FinalizeSubmission{
		submissionReader: submissionReader,
		moduleACL:        moduleACL,
		uow:              uow,
	}
}

func (s *FinalizeSubmission) Execute(ctx context.Context, command *FinalizeSubmissionCommand) (*response.FinalizeSubmissionResponse, error) {
	// Find submission by code
	submission, err := s.submissionReader.FindByCode(ctx, command.SubmissionCode)
	if err != nil {
		return nil, err
	}

	// Check if submission is in progress
	if !submission.IsInProgress() {
		return nil, constant.ErrSubmissionAlreadyDone
	}

	// Validate module exists and is published via ACL
	_, err = s.moduleACL.GetPublishedModule(ctx, command.ModuleSlug)
	if err != nil {
		return nil, err
	}

	// Finalize submission
	err = submission.Finalize()
	if err != nil {
		return nil, err
	}

	// Save via UnitOfWork
	tx, err := s.uow.Begin()
	if err != nil {
		return nil, err
	}

	err = tx.SubmissionWriter().Save(ctx, submission)
	if err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return nil, errRollback
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// Return finalization response
	return &response.FinalizeSubmissionResponse{
		StudentName: submission.StudentName,
		Score:       submission.Score(),
		Total:       submission.TotalQuestions,
		Status:      submission.Status.String(),
	}, nil
}
