package module

import (
	"context"
	"time"

	"github.com/arvinpaundra/private-api/core/util"
	"github.com/arvinpaundra/private-api/domain/module/entity"
	"github.com/arvinpaundra/private-api/domain/module/repository"
	"github.com/arvinpaundra/private-api/model"
	"github.com/guregu/null/v6"
	"gorm.io/gorm"
)

var _ repository.ModuleWriter = (*ModuleWriterRepository)(nil)

type ModuleWriterRepository struct {
	db *gorm.DB
}

func NewModuleWriterRepository(db *gorm.DB) *ModuleWriterRepository {
	return &ModuleWriterRepository{
		db: db,
	}
}

func (r *ModuleWriterRepository) Save(ctx context.Context, module *entity.Module) error {
	if module.IsUpdated() {
		return r.update(ctx, module)
	} else if module.IsRemoved() {
		return r.remove(ctx, module)
	}

	return r.insert(ctx, module)
}

func (r *ModuleWriterRepository) insert(ctx context.Context, module *entity.Module) error {
	moduleModel := model.Module{
		ID:          util.ParseUUID(module.ID),
		UserID:      util.ParseUUID(module.UserID),
		SubjectID:   util.ParseUUID(module.SubjectID),
		GradeID:     util.ParseUUID(module.GradeID),
		Title:       module.Title,
		Slug:        module.Slug,
		Description: null.StringFromPtr(module.Description),
		Type:        model.ModuleType(module.Type),
		IsPublished: module.IsPublished,
	}

	err := r.db.Model(&model.Module{}).WithContext(ctx).Create(&moduleModel).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *ModuleWriterRepository) update(ctx context.Context, module *entity.Module) error {
	// Update module fields using map to handle zero values
	updates := map[string]any{
		"subject_id":   util.ParseUUID(module.SubjectID),
		"grade_id":     util.ParseUUID(module.GradeID),
		"title":        module.Title,
		"slug":         module.Slug,
		"description":  null.StringFromPtr(module.Description),
		"type":         model.ModuleType(module.Type),
		"is_published": module.IsPublished,
	}

	err := r.db.Model(&model.Module{}).WithContext(ctx).Where("id = ?", module.ID).Updates(updates).Error
	if err != nil {
		return err
	}

	// Handle questions cascade
	for _, question := range module.Questions {
		if question.IsCreated() {
			// Insert new question
			questionModel := model.Question{
				ID:       util.ParseUUID(question.ID),
				ModuleID: util.ParseUUID(question.ModuleID),
				Content:  question.Content,
				Slug:     question.Slug,
			}

			err := r.db.Model(&model.Question{}).WithContext(ctx).Create(&questionModel).Error
			if err != nil {
				return err
			}

			// Insert question choices
			for _, choice := range question.Choices {
				choiceModel := model.QuestionChoice{
					ID:              util.ParseUUID(choice.ID),
					QuestionID:      util.ParseUUID(choice.QuestionID),
					Content:         choice.Content,
					IsCorrectAnswer: choice.IsCorrectAnswer,
				}

				err := r.db.Model(&model.QuestionChoice{}).WithContext(ctx).Create(&choiceModel).Error
				if err != nil {
					return err
				}
			}
		} else if question.IsUpdated() {
			// Update existing question using map to handle zero values
			updates := map[string]any{
				"content": question.Content,
				"slug":    question.Slug,
			}

			err := r.db.Model(&model.Question{}).WithContext(ctx).Where("id = ?", question.ID).Updates(updates).Error
			if err != nil {
				return err
			}

			// Handle question choices cascade
			for _, choice := range question.Choices {
				if choice.IsCreated() {
					// Insert new choice
					choiceModel := model.QuestionChoice{
						ID:              util.ParseUUID(choice.ID),
						QuestionID:      util.ParseUUID(choice.QuestionID),
						Content:         choice.Content,
						IsCorrectAnswer: choice.IsCorrectAnswer,
					}

					err := r.db.Model(&model.QuestionChoice{}).WithContext(ctx).Create(&choiceModel).Error
					if err != nil {
						return err
					}
				} else if choice.IsUpdated() {
					// Update existing choice using map to handle zero values
					updates := map[string]any{
						"content":           choice.Content,
						"is_correct_answer": choice.IsCorrectAnswer,
					}

					err := r.db.Model(&model.QuestionChoice{}).WithContext(ctx).Where("id = ?", choice.ID).Updates(updates).Error
					if err != nil {
						return err
					}
				} else if choice.IsRemoved() {
					// Soft delete choice
					choiceModel := model.QuestionChoice{
						DeletedAt: null.TimeFrom(time.Now().UTC()),
					}

					err := r.db.Model(&model.QuestionChoice{}).WithContext(ctx).Where("id = ?", choice.ID).Updates(&choiceModel).Error
					if err != nil {
						return err
					}
				}
			}
		} else if question.IsRemoved() {
			// Soft delete question
			questionModel := model.Question{
				DeletedAt: null.TimeFrom(time.Now().UTC()),
			}

			err := r.db.Model(&model.Question{}).WithContext(ctx).Where("id = ?", question.ID).Updates(&questionModel).Error
			if err != nil {
				return err
			}

			// Soft delete all associated choices
			err = r.db.Model(&model.QuestionChoice{}).WithContext(ctx).
				Where("question_id = ?", question.ID).
				Updates(model.QuestionChoice{DeletedAt: null.TimeFrom(time.Now().UTC())}).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *ModuleWriterRepository) remove(ctx context.Context, module *entity.Module) error {
	moduleModel := model.Module{
		DeletedAt: null.TimeFrom(time.Now().UTC()),
	}

	err := r.db.Model(&model.Module{}).WithContext(ctx).Where("id = ?", module.ID).Updates(&moduleModel).Error
	if err != nil {
		return err
	}

	return nil
}
