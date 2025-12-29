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
	moduleModel := model.Module{
		SubjectID:   util.ParseUUID(module.SubjectID),
		GradeID:     util.ParseUUID(module.GradeID),
		Title:       module.Title,
		Slug:        module.Slug,
		Description: null.StringFromPtr(module.Description),
		Type:        model.ModuleType(module.Type),
		IsPublished: module.IsPublished,
	}

	err := r.db.Model(&model.Module{}).WithContext(ctx).Where("id = ?", module.ID).Updates(&moduleModel).Error
	if err != nil {
		return err
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
