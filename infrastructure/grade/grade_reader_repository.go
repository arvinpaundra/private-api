package grade

import (
	"context"
	"errors"
	"strings"

	"github.com/arvinpaundra/private-api/domain/grade/constant"
	"github.com/arvinpaundra/private-api/domain/grade/entity"
	"github.com/arvinpaundra/private-api/domain/grade/repository"
	"github.com/arvinpaundra/private-api/model"
	"gorm.io/gorm"
)

var _ repository.GradeReader = (*GradeReaderRepository)(nil)

type GradeReaderRepository struct {
	db *gorm.DB
}

func NewGradeReaderRepository(db *gorm.DB) *GradeReaderRepository {
	return &GradeReaderRepository{
		db: db,
	}
}

func (r *GradeReaderRepository) HasSimilarGrade(ctx context.Context, name string, userID string) (bool, error) {
	var isExists bool

	err := r.db.WithContext(ctx).
		Raw(
			`SELECT EXISTS(SELECT 1 FROM grades WHERE LOWER(name) = ? AND user_id = ? AND deleted_at IS NULL)`,
			strings.ToLower(name),
			userID,
		).
		Scan(&isExists).Error

	if err != nil {
		return false, err
	}

	return isExists, nil
}

func (r *GradeReaderRepository) HasSimilarGradeExclusive(ctx context.Context, name string, userID string, excludeGradeID string) (bool, error) {
	var isExists bool

	err := r.db.WithContext(ctx).
		Raw(
			`SELECT EXISTS(SELECT 1 FROM grades WHERE LOWER(name) = ? AND user_id = ? AND id <> ? AND deleted_at IS NULL)`,
			strings.ToLower(name),
			userID,
			excludeGradeID,
		).
		Scan(&isExists).Error

	if err != nil {
		return false, err
	}

	return isExists, nil
}

func (r *GradeReaderRepository) IsGradeExist(ctx context.Context, gradeID string, userID string) (bool, error) {
	var isExists bool

	err := r.db.WithContext(ctx).
		Raw(
			`SELECT EXISTS(SELECT 1 FROM grades WHERE id = ? AND user_id = ? AND deleted_at IS NULL)`,
			gradeID,
			userID,
		).
		Scan(&isExists).Error

	if err != nil {
		return false, err
	}

	return isExists, nil
}

func (r *GradeReaderRepository) FindGradeByID(ctx context.Context, gradeID string, userID string) (*entity.Grade, error) {
	var gradeModel model.Grade

	err := r.db.Model(&model.Grade{}).
		WithContext(ctx).
		Select("id", "name", "description").
		Where("id = ?", gradeID).
		Where("user_id = ?", userID).
		Where("deleted_at IS NULL").
		First(&gradeModel).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constant.ErrGradeNotFound
		}

		return nil, err
	}

	grade := &entity.Grade{
		ID:          gradeModel.ID.String(),
		Name:        gradeModel.Name,
		Description: gradeModel.Description.Ptr(),
	}

	return grade, nil
}

func (r *GradeReaderRepository) AllGrades(ctx context.Context, userID string, keyword string) ([]*entity.Grade, error) {
	var gradeModels []model.Grade

	query := r.db.Model(&model.Grade{}).
		WithContext(ctx).
		Select("id", "name", "description").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Scopes(func(db *gorm.DB) *gorm.DB {
			if keyword != "" {
				return db.Where("name ILIKE ?", "%"+keyword+"%")
			}
			return db
		}).
		Find(&gradeModels).
		Error

	if query != nil {
		return nil, query
	}

	grades := make([]*entity.Grade, len(gradeModels))

	for i, gradeModel := range gradeModels {
		grades[i] = &entity.Grade{
			ID:          gradeModel.ID.String(),
			Name:        gradeModel.Name,
			Description: gradeModel.Description.Ptr(),
		}
	}

	return grades, nil
}
