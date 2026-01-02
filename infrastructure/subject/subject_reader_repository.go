package subject

import (
	"context"
	"errors"
	"strings"

	"github.com/arvinpaundra/private-api/domain/subject/constant"
	"github.com/arvinpaundra/private-api/domain/subject/entity"
	"github.com/arvinpaundra/private-api/domain/subject/repository"
	"github.com/arvinpaundra/private-api/model"
	"gorm.io/gorm"
)

var _ repository.SubjectReader = (*SubjectReaderRepository)(nil)

type SubjectReaderRepository struct {
	db *gorm.DB
}

func NewSubjectReaderRepository(db *gorm.DB) *SubjectReaderRepository {
	return &SubjectReaderRepository{
		db: db,
	}
}

func (r *SubjectReaderRepository) HasSimilarSubject(ctx context.Context, name string, userID string) (bool, error) {
	var isExists bool

	err := r.db.WithContext(ctx).
		Raw(
			`SELECT EXISTS(SELECT 1 FROM subjects WHERE LOWER(name) = ? AND user_id = ? AND deleted_at IS NULL)`,
			strings.ToLower(name),
			userID,
		).
		Scan(&isExists).Error

	if err != nil {
		return false, err
	}

	return isExists, nil
}

func (r *SubjectReaderRepository) HasSimilarSubjectExclusive(ctx context.Context, name string, userID string, excludeSubjectID string) (bool, error) {
	var isExists bool

	err := r.db.WithContext(ctx).
		Raw(
			`SELECT EXISTS(SELECT 1 FROM subjects WHERE LOWER(name) = ? AND user_id = ? AND id <> ? AND deleted_at IS NULL)`,
			strings.ToLower(name),
			userID,
			excludeSubjectID,
		).
		Scan(&isExists).Error

	if err != nil {
		return false, err
	}

	return isExists, nil
}

func (r *SubjectReaderRepository) IsSubjectExist(ctx context.Context, subjectID string, userID string) (bool, error) {
	var isExists bool

	err := r.db.WithContext(ctx).
		Raw(
			`SELECT EXISTS(SELECT 1 FROM subjects WHERE id = ? AND user_id = ? AND deleted_at IS NULL)`,
			subjectID,
			userID,
		).
		Scan(&isExists).Error

	if err != nil {
		return false, err
	}

	return isExists, nil
}

func (r *SubjectReaderRepository) FindSubjectByID(ctx context.Context, subjectID string, userID string) (*entity.Subject, error) {
	var subjectModel model.Subject

	err := r.db.Model(&model.Subject{}).
		WithContext(ctx).
		Select("id", "name", "description").
		Where("id = ?", subjectID).
		Where("user_id = ?", userID).
		Where("deleted_at IS NULL").
		First(&subjectModel).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constant.ErrSubjectNotFound
		}

		return nil, err
	}

	subject := &entity.Subject{
		ID:          subjectModel.ID.String(),
		Name:        subjectModel.Name,
		Description: subjectModel.Description.Ptr(),
	}

	return subject, nil
}

func (r *SubjectReaderRepository) AllSubjects(ctx context.Context, userID string, keyword string) ([]*entity.Subject, error) {
	var subjectModels []model.Subject

	query := r.db.Model(&model.Subject{}).
		WithContext(ctx).
		Select("id", "name", "description").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if keyword != "" {
				return db.Where("name ILIKE ?", "%"+keyword+"%")
			}
			return db
		}).
		Find(&subjectModels).
		Error

	if query != nil {
		return nil, query
	}

	subjects := make([]*entity.Subject, len(subjectModels))

	for i, subjectModel := range subjectModels {
		subjects[i] = &entity.Subject{
			ID:          subjectModel.ID.String(),
			Name:        subjectModel.Name,
			Description: subjectModel.Description.Ptr(),
		}
	}

	return subjects, nil
}

func (r *SubjectReaderRepository) CountByUserID(ctx context.Context, userID string) (int, error) {
	var count int64

	err := r.db.Model(&model.Subject{}).
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
