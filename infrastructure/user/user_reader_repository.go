package user

import (
	"context"
	"errors"

	"github.com/arvinpaundra/private-api/domain/user/constant"
	"github.com/arvinpaundra/private-api/domain/user/entity"
	"github.com/arvinpaundra/private-api/domain/user/repository"
	"github.com/arvinpaundra/private-api/model"
	"gorm.io/gorm"
)

var _ repository.UserReader = (*UserReaderRepository)(nil)

type UserReaderRepository struct {
	db *gorm.DB
}

func NewUserReaderRepository(db *gorm.DB) *UserReaderRepository {
	return &UserReaderRepository{db: db}
}

func (r *UserReaderRepository) FindByID(ctx context.Context, userID string) (*entity.User, error) {
	var userModel model.User

	err := r.db.Model(&model.User{}).
		WithContext(ctx).
		Where("id = ?", userID).
		First(&userModel).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constant.ErrUserNotFound
		}

		return nil, err
	}

	user := entity.User{
		ID:       userModel.ID.String(),
		Email:    userModel.Email,
		Fullname: userModel.Fullname,
	}

	return &user, nil
}
