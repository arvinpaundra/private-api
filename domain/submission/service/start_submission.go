package service

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/submission/entity"
	"github.com/arvinpaundra/private-api/domain/submission/repository"
	"github.com/arvinpaundra/private-api/domain/submission/response"
)

type StartSubmissionCommand struct {
	ModuleSlug  string `json:"-" validate:"required"`
	StudentName string `json:"student_name" validate:"required,min=3,max=100"`
}

type StartSubmission struct {
	uow       repository.UnitOfWork
	moduleACL repository.ModuleACL
}

func NewStartSubmission(
	uow repository.UnitOfWork,
	moduleACL repository.ModuleACL,
) *StartSubmission {
	return &StartSubmission{
		uow:       uow,
		moduleACL: moduleACL,
	}
}

func (s *StartSubmission) Execute(ctx context.Context, command *StartSubmissionCommand) (*response.StartSubmissionResponse, error) {
	// Validate module exists and is published via ACL
	module, err := s.moduleACL.GetPublishedModule(ctx, command.ModuleSlug)
	if err != nil {
		return nil, err
	}

	// Get total questions count
	totalQuestions, err := s.moduleACL.GetTotalQuestions(ctx, module.Slug)
	if err != nil {
		return nil, err
	}

	// Create new submission with generated code
	submission, err := entity.NewSubmission(module.ID, command.StudentName)
	if err != nil {
		return nil, err
	}

	submission.SetTotalQuestions(totalQuestions)

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

	// Return submission code and status
	return &response.StartSubmissionResponse{
		Code:   submission.Code,
		Status: submission.Status.String(),
	}, nil
}
