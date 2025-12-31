package submission

import (
	"context"
	"time"

	"github.com/arvinpaundra/private-api/core/util"
	"github.com/arvinpaundra/private-api/domain/submission/entity"
	"github.com/arvinpaundra/private-api/domain/submission/repository"
	"github.com/arvinpaundra/private-api/model"
	"github.com/guregu/null/v6"
	"gorm.io/gorm"
)

var _ repository.SubmissionWriter = (*SubmissionWriterRepository)(nil)

type SubmissionWriterRepository struct {
	db *gorm.DB
}

func NewSubmissionWriterRepository(db *gorm.DB) *SubmissionWriterRepository {
	return &SubmissionWriterRepository{
		db: db,
	}
}

func (r *SubmissionWriterRepository) Save(ctx context.Context, submission *entity.Submission) error {
	if submission.IsUpdated() {
		return r.update(ctx, submission)
	}

	return r.insert(ctx, submission)
}

func (r *SubmissionWriterRepository) insert(ctx context.Context, submission *entity.Submission) error {
	submissionModel := model.Submission{
		ID:             util.ParseUUID(submission.ID),
		ModuleID:       util.ParseUUID(submission.ModuleID),
		Code:           submission.Code,
		StudentName:    submission.StudentName,
		Status:         model.SubmissionStatus(submission.Status),
		TotalQuestions: submission.TotalQuestions,
		SubmittedAt:    null.TimeFromPtr(submission.SubmittedAt),
	}

	err := r.db.Model(&model.Submission{}).WithContext(ctx).Create(&submissionModel).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *SubmissionWriterRepository) update(ctx context.Context, submission *entity.Submission) error {
	// Update submission fields using map to handle zero values
	updates := map[string]any{
		"student_name":    submission.StudentName,
		"status":          model.SubmissionStatus(submission.Status),
		"total_questions": submission.TotalQuestions,
		"submitted_at":    null.TimeFromPtr(submission.SubmittedAt),
	}

	err := r.db.Model(&model.Submission{}).
		WithContext(ctx).
		Where("id = ?", submission.ID).
		Updates(updates).
		Error
	if err != nil {
		return err
	}

	// Handle submission answers cascade
	for _, answer := range submission.Answers {
		if answer.IsCreated() {
			// Insert new answer
			answerModel := model.SubmissionAnswer{
				ID:           util.ParseUUID(answer.ID),
				SubmissionID: util.ParseUUID(answer.SubmissionID),
				QuestionSlug: answer.QuestionSlug,
				Question:     answer.Question,
				Answer:       answer.Answer,
				IsCorrect:    answer.IsCorrect,
			}

			err := r.db.Model(&model.SubmissionAnswer{}).WithContext(ctx).Create(&answerModel).Error
			if err != nil {
				return err
			}
		} else if answer.IsUpdated() {
			// Update existing answer
			answerUpdates := map[string]any{
				"question":   answer.Question,
				"answer":     answer.Answer,
				"is_correct": answer.IsCorrect,
			}

			err := r.db.Model(&model.SubmissionAnswer{}).
				WithContext(ctx).
				Where("id = ?", answer.ID).
				Updates(answerUpdates).
				Error
			if err != nil {
				return err
			}
		} else if answer.IsRemoved() {
			// Soft delete answer
			now := time.Now().UTC()
			err := r.db.Model(&model.SubmissionAnswer{}).
				WithContext(ctx).
				Where("id = ?", answer.ID).
				Update("deleted_at", now).
				Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}
