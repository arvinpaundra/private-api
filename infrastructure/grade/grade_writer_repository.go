package grade

import (
	"context"
	"time"

	"github.com/arvinpaundra/private-api/core/util"
	"github.com/arvinpaundra/private-api/domain/grade/entity"
	"github.com/arvinpaundra/private-api/domain/grade/repository"
	"github.com/arvinpaundra/private-api/model"
	"github.com/guregu/null/v6"
	"gorm.io/gorm"
)

var _ repository.GradeWriter = (*GradeWriterRepository)(nil)

type GradeWriterRepository struct {
	db *gorm.DB
}

func NewGradeWriterRepository(db *gorm.DB) *GradeWriterRepository {
	return &GradeWriterRepository{
		db: db,
	}
}

func (r *GradeWriterRepository) Save(ctx context.Context, grade *entity.Grade) error {
	if grade.IsUpdated() {
		return r.update(ctx, grade)
	} else if grade.IsRemoved() {
		return r.remove(ctx, grade)
	}

	return r.insert(ctx, grade)
}

func (r *GradeWriterRepository) insert(ctx context.Context, grade *entity.Grade) error {
	gradeModel := model.Grade{
		ID:          util.ParseUUID(grade.ID),
		UserID:      util.ParseUUID(grade.UserID),
		Name:        grade.Name,
		Description: null.StringFromPtr(grade.Description),
	}

	err := r.db.Model(&model.Grade{}).WithContext(ctx).Create(&gradeModel).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *GradeWriterRepository) update(ctx context.Context, grade *entity.Grade) error {
	gradeModel := model.Grade{
		Name:        grade.Name,
		Description: null.StringFromPtr(grade.Description),
	}

	err := r.db.Model(&model.Grade{}).WithContext(ctx).Where("id = ?", grade.ID).Updates(&gradeModel).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *GradeWriterRepository) remove(ctx context.Context, grade *entity.Grade) error {
	gradeModel := model.Grade{
		DeletedAt: null.TimeFrom(time.Now().UTC()),
	}

	err := r.db.Model(&model.Grade{}).WithContext(ctx).Where("id = ?", grade.ID).Updates(&gradeModel).Error
	if err != nil {
		return err
	}

	return nil
}
