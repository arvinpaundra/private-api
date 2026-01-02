package module

import (
	"context"
	"errors"

	"github.com/arvinpaundra/private-api/domain/module/constant"
	"github.com/arvinpaundra/private-api/domain/module/entity"
	"github.com/arvinpaundra/private-api/domain/module/repository"
	"github.com/arvinpaundra/private-api/model"
	"gorm.io/gorm"
)

var _ repository.ModuleReader = (*ModuleReaderRepository)(nil)

type ModuleReaderRepository struct {
	db *gorm.DB
}

func NewModuleReaderRepository(db *gorm.DB) *ModuleReaderRepository {
	return &ModuleReaderRepository{
		db: db,
	}
}

func (r *ModuleReaderRepository) FindByID(ctx context.Context, moduleID, userID string) (*entity.Module, error) {
	var module model.Module

	err := r.db.Model(&model.Module{}).
		WithContext(ctx).
		Select("id", "user_id", "subject_id", "grade_id", "title", "slug", "description", "type", "is_published").
		Where("id = ?", moduleID).
		Where("user_id = ?", userID).
		Where("deleted_at IS NULL").
		First(&module).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constant.ErrModuleNotFound
		}
		return nil, err
	}

	return &entity.Module{
		ID:          module.ID.String(),
		UserID:      module.UserID.String(),
		SubjectID:   module.SubjectID.String(),
		GradeID:     module.GradeID.String(),
		Title:       module.Title,
		Slug:        module.Slug,
		Description: module.Description.Ptr(),
		Type:        constant.ModuleType(module.Type),
		IsPublished: module.IsPublished,
	}, nil
}

func (r *ModuleReaderRepository) FindBySlug(ctx context.Context, slug, userID string) (*entity.Module, error) {
	var module model.Module

	err := r.db.Model(&model.Module{}).
		WithContext(ctx).
		Select("id", "user_id", "subject_id", "grade_id", "title", "slug", "description", "type", "is_published").
		Where("slug = ?", slug).
		Where("user_id = ?", userID).
		Where("deleted_at IS NULL").
		First(&module).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constant.ErrModuleNotFound
		}
		return nil, err
	}

	return &entity.Module{
		ID:          module.ID.String(),
		UserID:      module.UserID.String(),
		SubjectID:   module.SubjectID.String(),
		GradeID:     module.GradeID.String(),
		Title:       module.Title,
		Slug:        module.Slug,
		Description: module.Description.Ptr(),
		Type:        constant.ModuleType(module.Type),
		IsPublished: module.IsPublished,
	}, nil
}

func (r *ModuleReaderRepository) FindModuleDetailBySlug(ctx context.Context, slug, userID string) (*entity.Module, error) {
	var module model.Module

	err := r.db.Model(&model.Module{}).
		WithContext(ctx).
		Select("id", "user_id", "subject_id", "grade_id", "title", "slug", "description", "type", "is_published").
		Where("slug = ?", slug).
		Where("user_id = ?", userID).
		Where("deleted_at IS NULL").
		Preload("Questions", "deleted_at IS NULL").
		Preload("Questions.Choices", "deleted_at IS NULL").
		First(&module).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constant.ErrModuleNotFound
		}
		return nil, err
	}

	// Convert questions with choices
	questions := make([]*entity.Question, len(module.Questions))

	for i, question := range module.Questions {
		choices := make([]*entity.QuestionChoice, len(question.Choices))

		for j, choice := range question.Choices {
			choices[j] = &entity.QuestionChoice{
				ID:              choice.ID.String(),
				QuestionID:      choice.QuestionID.String(),
				Content:         choice.Content,
				IsCorrectAnswer: choice.IsCorrectAnswer,
			}
		}

		questions[i] = &entity.Question{
			ID:       question.ID.String(),
			ModuleID: question.ModuleID.String(),
			Content:  question.Content,
			Slug:     question.Slug,
			Choices:  choices,
		}
	}

	return &entity.Module{
		ID:          module.ID.String(),
		UserID:      module.UserID.String(),
		SubjectID:   module.SubjectID.String(),
		GradeID:     module.GradeID.String(),
		Title:       module.Title,
		Slug:        module.Slug,
		Description: module.Description.Ptr(),
		Type:        constant.ModuleType(module.Type),
		IsPublished: module.IsPublished,
		Questions:   questions,
	}, nil
}

func (r *ModuleReaderRepository) TotalModules(ctx context.Context, userID, subjectID, gradeID, keyword string) (int, error) {
	var total int64

	err := r.db.Model(&model.Module{}).
		WithContext(ctx).
		Select("id").
		Where("user_id = ?", userID).
		Where("deleted_at IS NULL").
		Scopes(func(db *gorm.DB) *gorm.DB {
			if keyword != "" {
				return db.Where("title ILIKE ?", "%"+keyword+"%")
			}
			return db
		}).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if subjectID != "" {
				return db.Where("subject_id = ?", subjectID)
			}
			return db
		}).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if gradeID != "" {
				return db.Where("grade_id = ?", gradeID)
			}
			return db
		}).
		Count(&total).
		Error

	if err != nil {
		return 0, err
	}

	return int(total), nil
}

func (r *ModuleReaderRepository) FindAllModules(ctx context.Context, userID, subjectID, gradeID, keyword string, limit, offset int) ([]*entity.Module, error) {
	type ModuleWithCount struct {
		model.Module
		QuestionsCount int `gorm:"column:questions_count"`
	}

	var modulesWithCount []ModuleWithCount

	// Get paginated results with question count
	err := r.db.Model(&model.Module{}).
		WithContext(ctx).
		Where("modules.user_id = ?", userID).
		Where("modules.deleted_at IS NULL").
		Scopes(func(db *gorm.DB) *gorm.DB {
			if keyword != "" {
				return db.Where("modules.title ILIKE ?", "%"+keyword+"%")
			}
			return db
		}).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if subjectID != "" {
				return db.Where("modules.subject_id = ?", subjectID)
			}
			return db
		}).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if gradeID != "" {
				return db.Where("modules.grade_id = ?", gradeID)
			}
			return db
		}).
		Select("modules.id", "modules.user_id", "modules.subject_id", "modules.grade_id",
			"modules.title", "modules.description", "modules.is_published", "modules.slug", "modules.type",
			"COUNT(questions.id) as questions_count").
		Joins("LEFT JOIN questions ON questions.module_id = modules.id AND questions.deleted_at IS NULL").
		Group("modules.id").
		Order("modules.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&modulesWithCount).
		Error

	if err != nil {
		return nil, err
	}

	// Convert to domain entities
	modules := make([]*entity.Module, len(modulesWithCount))
	for i, m := range modulesWithCount {
		// Create question placeholders to represent count
		questions := make([]*entity.Question, m.QuestionsCount)

		modules[i] = &entity.Module{
			ID:          m.Module.ID.String(),
			UserID:      m.Module.UserID.String(),
			SubjectID:   m.Module.SubjectID.String(),
			GradeID:     m.Module.GradeID.String(),
			Title:       m.Module.Title,
			Slug:        m.Slug,
			Description: m.Module.Description.Ptr(),
			Type:        constant.ModuleType(m.Module.Type),
			IsPublished: m.Module.IsPublished,
			Questions:   questions,
		}
	}

	return modules, nil
}

func (r *ModuleReaderRepository) FindPublishedModuleBySlug(ctx context.Context, slug string) (*entity.Module, error) {
	var module model.Module

	err := r.db.Model(&model.Module{}).
		WithContext(ctx).
		Select("id", "user_id", "subject_id", "grade_id", "title", "slug", "description", "type", "is_published").
		Where("slug = ?", slug).
		Where("is_published = true").
		Where("deleted_at IS NULL").
		First(&module).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constant.ErrModuleNotFound
		}
		return nil, err
	}

	return &entity.Module{
		ID:          module.ID.String(),
		UserID:      module.UserID.String(),
		SubjectID:   module.SubjectID.String(),
		GradeID:     module.GradeID.String(),
		Title:       module.Title,
		Slug:        module.Slug,
		Description: module.Description.Ptr(),
		Type:        constant.ModuleType(module.Type),
		IsPublished: module.IsPublished,
	}, nil
}

func (r *ModuleReaderRepository) FindPublishedQuestionBySlug(ctx context.Context, moduleSlug, questionSlug string) (*entity.Question, error) {
	var question model.Question

	err := r.db.Model(&model.Question{}).
		WithContext(ctx).
		Select("questions.id", "questions.module_id", "questions.content", "questions.slug").
		Joins("JOIN modules ON modules.id = questions.module_id").
		Where("modules.slug = ?", moduleSlug).
		Where("modules.is_published = true").
		Where("modules.deleted_at IS NULL").
		Where("questions.slug = ?", questionSlug).
		Where("questions.deleted_at IS NULL").
		Preload("Choices", "deleted_at IS NULL").
		First(&question).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constant.ErrQuestionNotFound
		}
		return nil, err
	}

	// Convert choices to domain entities
	choices := make([]*entity.QuestionChoice, len(question.Choices))

	for i, choice := range question.Choices {
		choices[i] = &entity.QuestionChoice{
			ID:              choice.ID.String(),
			QuestionID:      choice.QuestionID.String(),
			Content:         choice.Content,
			IsCorrectAnswer: choice.IsCorrectAnswer,
		}
	}

	return &entity.Question{
		ID:       question.ID.String(),
		ModuleID: question.ModuleID.String(),
		Content:  question.Content,
		Slug:     question.Slug,
		Choices:  choices,
	}, nil
}

func (r *ModuleReaderRepository) FindNextPublishedQuestion(ctx context.Context, moduleSlug, currentQuestionSlug string) (*entity.Question, error) {
	var nextQuestion model.Question

	err := r.db.Model(&model.Question{}).
		WithContext(ctx).
		Select("questions.id", "questions.module_id", "questions.content", "questions.slug").
		Joins("JOIN modules ON modules.id = questions.module_id").
		Where("modules.slug = ?", moduleSlug).
		Where("modules.is_published = true").
		Where("modules.deleted_at IS NULL").
		Where("questions.deleted_at IS NULL").
		Where("questions.created_at > (SELECT created_at FROM questions WHERE slug = ? AND deleted_at IS NULL)", currentQuestionSlug).
		Order("questions.created_at ASC").
		First(&nextQuestion).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &entity.Question{
		ID:       nextQuestion.ID.String(),
		ModuleID: nextQuestion.ModuleID.String(),
		Content:  nextQuestion.Content,
		Slug:     nextQuestion.Slug,
	}, nil
}

func (r *ModuleReaderRepository) CountQuestionsByModuleSlug(ctx context.Context, moduleSlug string) (int, error) {
	var count int64

	query := r.db.Model(&model.Question{}).
		WithContext(ctx).
		Joins("JOIN modules ON questions.module_id = modules.id").
		Where("modules.slug = ?", moduleSlug).
		Where("modules.is_published = true").
		Where("modules.deleted_at IS NULL").
		Where("questions.deleted_at IS NULL").
		Count(&count).
		Error

	if query != nil {
		return 0, query
	}

	return int(count), nil
}

func (r *ModuleReaderRepository) CountByUserID(ctx context.Context, userID string) (int, error) {
	var count int64

	err := r.db.Model(&model.Module{}).
		WithContext(ctx).
		Where("user_id = ?", userID).
		Where("deleted_at IS NULL").
		Count(&count).
		Error

	if err != nil {
		return 0, err
	}

	return int(count), nil
}
