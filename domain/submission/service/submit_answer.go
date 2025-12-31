package service

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/submission/constant"
	"github.com/arvinpaundra/private-api/domain/submission/entity"
	"github.com/arvinpaundra/private-api/domain/submission/repository"
	"github.com/arvinpaundra/private-api/domain/submission/response"
)

type SubmitAnswerCommand struct {
	SubmissionCode string `json:"-" validate:"required"`
	ModuleSlug     string `json:"-" validate:"required"`
	QuestionSlug   string `json:"question_slug" validate:"required"`
	ChoiceID       string `json:"choice_id" validate:"required"`
}

type SubmitAnswer struct {
	submissionReader repository.SubmissionReader
	uow              repository.UnitOfWork
	moduleACL        repository.ModuleACL
}

func NewSubmitAnswer(
	submissionReader repository.SubmissionReader,
	uow repository.UnitOfWork,
	moduleACL repository.ModuleACL,
) *SubmitAnswer {
	return &SubmitAnswer{
		submissionReader: submissionReader,
		uow:              uow,
		moduleACL:        moduleACL,
	}
}

func (s *SubmitAnswer) Execute(ctx context.Context, command *SubmitAnswerCommand) (*response.SubmitAnswerResponse, error) {
	// Find submission by code
	submission, err := s.submissionReader.FindByCode(ctx, command.SubmissionCode)
	if err != nil {
		return nil, err
	}

	// Check if submission is in progress
	if !submission.IsInProgress() {
		return nil, constant.ErrSubmissionAlreadyDone
	}

	// Prevent duplicate answers for the same question
	if submission.HasAnsweredQuestion(command.QuestionSlug) {
		return nil, constant.ErrDuplicateAnswer
	}

	// Validate module exists and is published
	module, err := s.moduleACL.GetPublishedModule(ctx, command.ModuleSlug)
	if err != nil {
		return nil, err
	}

	// Get question details to store question content
	question, err := s.moduleACL.GetQuestionBySlug(ctx, module.Slug, command.QuestionSlug)
	if err != nil {
		return nil, err
	}

	// Find the submitted choice to get its content
	submittedChoice, exist := question.GetChoiceByID(command.ChoiceID)
	if !exist {
		return nil, constant.ErrChoiceNotFound
	}

	// Get correct answer from module domain
	correctChoice, err := s.moduleACL.GetCorrectAnswer(ctx, module.Slug, question.Slug)
	if err != nil {
		return nil, err
	}

	// Determine if answer is correct
	isCorrect := command.ChoiceID == correctChoice.ID

	// Create submission answer with question content and choice content
	answer := entity.NewSubmissionAnswer(submission.ID, question.Slug, question.Content, submittedChoice.Content, isCorrect)

	// Add answer to submission
	err = submission.AddAnswer(answer)
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

	return &response.SubmitAnswerResponse{
		IsCorrect:            isCorrect,
		CorrectChoiceID:      correctChoice.ID,
		CorrectChoiceContent: correctChoice.Content,
		NextQuestionSlug:     question.NextQuestionSlug,
	}, nil
}
