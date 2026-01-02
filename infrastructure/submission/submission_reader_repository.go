package submission

import (
	"context"
	"time"

	"github.com/arvinpaundra/private-api/domain/submission/constant"
	"github.com/arvinpaundra/private-api/domain/submission/entity"
	"github.com/arvinpaundra/private-api/domain/submission/repository"
	"github.com/arvinpaundra/private-api/model"
	"gorm.io/gorm"
)

var _ repository.SubmissionReader = (*SubmissionReaderRepository)(nil)

type SubmissionReaderRepository struct {
	db *gorm.DB
}

func NewSubmissionReaderRepository(db *gorm.DB) *SubmissionReaderRepository {
	return &SubmissionReaderRepository{
		db: db,
	}
}

func (r *SubmissionReaderRepository) FindByCode(ctx context.Context, code string) (*entity.Submission, error) {
	var submissionModel model.Submission

	err := r.db.Model(&model.Submission{}).
		WithContext(ctx).
		Preload("Answers").
		Where("code = ?", code).
		First(&submissionModel).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, constant.ErrSubmissionNotFound
		}
		return nil, err
	}

	var submittedAt *time.Time
	if submissionModel.SubmittedAt.Valid {
		submittedAt = &submissionModel.SubmittedAt.Time
	}

	submission := &entity.Submission{
		ID:             submissionModel.ID.String(),
		ModuleID:       submissionModel.ModuleID.String(),
		Code:           submissionModel.Code,
		StudentName:    submissionModel.StudentName,
		Status:         constant.SubmissionStatus(submissionModel.Status),
		TotalQuestions: submissionModel.TotalQuestions,
		SubmittedAt:    submittedAt,
		Answers:        make([]*entity.SubmissionAnswer, len(submissionModel.Answers)),
	}

	for i, answerModel := range submissionModel.Answers {
		submission.Answers[i] = &entity.SubmissionAnswer{
			ID:           answerModel.ID.String(),
			SubmissionID: answerModel.SubmissionID.String(),
			QuestionSlug: answerModel.QuestionSlug,
			Question:     answerModel.Question,
			Answer:       answerModel.Answer,
			IsCorrect:    answerModel.IsCorrect,
		}
	}

	return submission, nil
}

func (r *SubmissionReaderRepository) TotalSubmissions(ctx context.Context, moduleID, status, keyword string) (int, error) {
	var count int64

	query := r.db.Model(&model.Submission{}).
		WithContext(ctx).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if moduleID != "" {
				db.Where("module_id = ?", moduleID)
			}
			if status != "" {
				db.Where("status = ?", status)
			}
			if keyword != "" {
				db.Where("student_name ILIKE ?", "%"+keyword+"%")
			}
			return db
		}).
		Count(&count).
		Error

	if query != nil {
		return 0, query
	}

	return int(count), nil
}

func (r *SubmissionReaderRepository) FindAllSubmissions(ctx context.Context, moduleID, status, keyword string, limit, offset int) ([]*entity.Submission, error) {
	var submissionModels []model.Submission

	query := r.db.Model(&model.Submission{}).
		WithContext(ctx).
		Preload("Answers").
		Scopes(func(db *gorm.DB) *gorm.DB {
			if moduleID != "" {
				db.Where("module_id = ?", moduleID)
			}
			if status != "" {
				db.Where("status = ?", status)
			}
			if keyword != "" {
				db.Where("student_name ILIKE ?", "%"+keyword+"%")
			}
			if limit > 0 {
				db.Limit(limit)
			}
			if offset > 0 {
				db.Offset(offset)
			}
			return db
		}).
		Order("created_at DESC").
		Find(&submissionModels).
		Error

	if query != nil {
		return nil, query
	}

	submissions := make([]*entity.Submission, len(submissionModels))

	for i, submissionModel := range submissionModels {
		var submittedAt *time.Time
		if submissionModel.SubmittedAt.Valid {
			submittedAt = &submissionModel.SubmittedAt.Time
		}

		submission := &entity.Submission{
			ID:             submissionModel.ID.String(),
			ModuleID:       submissionModel.ModuleID.String(),
			Code:           submissionModel.Code,
			StudentName:    submissionModel.StudentName,
			Status:         constant.SubmissionStatus(submissionModel.Status),
			TotalQuestions: submissionModel.TotalQuestions,
			SubmittedAt:    submittedAt,
			Answers:        make([]*entity.SubmissionAnswer, len(submissionModel.Answers)),
		}

		for j, answerModel := range submissionModel.Answers {
			submission.Answers[j] = &entity.SubmissionAnswer{
				ID:           answerModel.ID.String(),
				SubmissionID: answerModel.SubmissionID.String(),
				QuestionSlug: answerModel.QuestionSlug,
				Question:     answerModel.Question,
				Answer:       answerModel.Answer,
				IsCorrect:    answerModel.IsCorrect,
			}
		}

		submissions[i] = submission
	}

	return submissions, nil
}

func (r *SubmissionReaderRepository) FindAllSubmittedGroupedByModule(ctx context.Context, moduleIDs []string) (map[string][]*entity.Submission, error) {
	if len(moduleIDs) == 0 {
		return make(map[string][]*entity.Submission), nil
	}

	var submissionModels []model.Submission

	err := r.db.Model(&model.Submission{}).
		WithContext(ctx).
		Preload("Answers").
		Where("module_id IN ?", moduleIDs).
		Where("status = ?", constant.Submitted).
		Order("module_id, submitted_at DESC").
		Find(&submissionModels).
		Error

	if err != nil {
		return nil, err
	}

	// Group submissions by module ID
	grouped := make(map[string][]*entity.Submission)

	for _, submissionModel := range submissionModels {
		var submittedAt *time.Time
		if submissionModel.SubmittedAt.Valid {
			submittedAt = &submissionModel.SubmittedAt.Time
		}

		submission := &entity.Submission{
			ID:             submissionModel.ID.String(),
			ModuleID:       submissionModel.ModuleID.String(),
			Code:           submissionModel.Code,
			StudentName:    submissionModel.StudentName,
			Status:         constant.SubmissionStatus(submissionModel.Status),
			TotalQuestions: submissionModel.TotalQuestions,
			SubmittedAt:    submittedAt,
			Answers:        make([]*entity.SubmissionAnswer, len(submissionModel.Answers)),
		}

		// Map answers for score calculation
		for j, answerModel := range submissionModel.Answers {
			submission.Answers[j] = &entity.SubmissionAnswer{
				ID:           answerModel.ID.String(),
				SubmissionID: answerModel.SubmissionID.String(),
				QuestionSlug: answerModel.QuestionSlug,
				Question:     answerModel.Question,
				Answer:       answerModel.Answer,
				IsCorrect:    answerModel.IsCorrect,
			}
		}

		moduleID := submissionModel.ModuleID.String()
		grouped[moduleID] = append(grouped[moduleID], submission)
	}

	return grouped, nil
}

func (r *SubmissionReaderRepository) CountSubmitted(ctx context.Context) (int, error) {
	var count int64

	err := r.db.Model(&model.Submission{}).
		WithContext(ctx).
		Where("status = ?", model.Submitted).
		Count(&count).
		Error

	if err != nil {
		return 0, err
	}

	return int(count), nil
}
