package entity

import (
	"time"

	"github.com/arvinpaundra/private-api/core/trait"
	"github.com/arvinpaundra/private-api/core/util"
	"github.com/arvinpaundra/private-api/domain/submission/constant"
)

type Submission struct {
	trait.Createable
	trait.Updateable
	trait.Removeable

	ID             string
	ModuleID       string
	Code           string
	StudentName    string
	Status         constant.SubmissionStatus
	TotalQuestions int
	SubmittedAt    *time.Time

	Answers []*SubmissionAnswer
}

func NewSubmission(moduleID, studentName string) (*Submission, error) {
	code, err := util.RandomAlphanumeric(16)
	if err != nil {
		return nil, err
	}

	submission := &Submission{
		ID:          util.GenerateUUID(),
		ModuleID:    moduleID,
		Code:        code,
		StudentName: studentName,
		Status:      constant.InProgress,
	}

	submission.MarkCreate()

	return submission, nil
}

func (s *Submission) Submit() error {
	if !s.IsInProgress() {
		return constant.ErrCannotSubmit
	}

	s.Status = constant.Submitted
	s.MarkUpdate()

	return nil
}

func (s *Submission) Cancel() error {
	if !s.IsInProgress() {
		return constant.ErrCannotCancel
	}

	s.Status = constant.Canceled
	s.MarkUpdate()

	return nil
}

func (s *Submission) AddAnswer(answer *SubmissionAnswer) error {
	if !s.IsInProgress() {
		return constant.ErrSubmissionAlreadyDone
	}

	s.Answers = append(s.Answers, answer)
	s.MarkUpdate()

	return nil
}

func (s *Submission) IsInProgress() bool {
	return s.Status == constant.InProgress
}

func (s *Submission) IsSubmitted() bool {
	return s.Status == constant.Submitted
}

func (s *Submission) IsCanceled() bool {
	return s.Status == constant.Canceled
}

func (s *Submission) HasAnsweredQuestion(questionSlug string) bool {
	for _, answer := range s.Answers {
		if answer.QuestionSlug == questionSlug {
			return true
		}
	}
	return false
}

func (s *Submission) Finalize() error {
	if err := s.Submit(); err != nil {
		return err
	}

	submittedAt := time.Now().UTC()

	s.SubmittedAt = &submittedAt

	return nil
}

func (s *Submission) Score() int {
	score := 0

	for _, answer := range s.Answers {
		if answer.IsCorrect {
			score++
		}
	}

	return score
}

func (s *Submission) SetTotalQuestions(total int) {
	s.TotalQuestions = total
}
