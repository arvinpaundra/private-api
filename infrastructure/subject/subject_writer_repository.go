package subject

import (
	"context"
	"time"

	"github.com/arvinpaundra/private-api/core/util"
	"github.com/arvinpaundra/private-api/domain/subject/entity"
	"github.com/arvinpaundra/private-api/domain/subject/repository"
	"github.com/arvinpaundra/private-api/model"
	"github.com/guregu/null/v6"
	"gorm.io/gorm"
)

var _ repository.SubjectWriter = (*SubjectWriterRepository)(nil)

type SubjectWriterRepository struct {
	db *gorm.DB
}

func NewSubjectWriterRepository(db *gorm.DB) *SubjectWriterRepository {
	return &SubjectWriterRepository{
		db: db,
	}
}

func (r *SubjectWriterRepository) Save(ctx context.Context, subject *entity.Subject) error {
	if subject.IsUpdated() {
		return r.update(ctx, subject)
	} else if subject.IsRemoved() {
		return r.remove(ctx, subject)
	}

	return r.insert(ctx, subject)
}

func (r *SubjectWriterRepository) insert(ctx context.Context, subject *entity.Subject) error {
	subjectModel := model.Subject{
		ID:          util.ParseUUID(subject.ID),
		UserID:      util.ParseUUID(subject.UserID),
		Name:        subject.Name,
		Description: null.StringFromPtr(subject.Description),
	}

	err := r.db.Model(&model.Subject{}).WithContext(ctx).Create(&subjectModel).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *SubjectWriterRepository) update(ctx context.Context, subject *entity.Subject) error {
	subjectModel := model.Subject{
		Name:        subject.Name,
		Description: null.StringFromPtr(subject.Description),
	}

	err := r.db.Model(&model.Subject{}).WithContext(ctx).Where("id = ?", subject.ID).Updates(&subjectModel).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *SubjectWriterRepository) remove(ctx context.Context, subject *entity.Subject) error {
	subjectModel := model.Subject{
		DeletedAt: null.TimeFrom(time.Now().UTC()),
	}

	err := r.db.Model(&model.Subject{}).WithContext(ctx).Where("id = ?", subject.ID).Updates(&subjectModel).Error
	if err != nil {
		return err
	}

	return nil
}
