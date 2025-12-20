package auth

import (
	"context"
	"errors"

	"github.com/arvinpaundra/private-api/domain/auth/constant"
	"github.com/arvinpaundra/private-api/domain/auth/entity"
	"github.com/arvinpaundra/private-api/domain/auth/repository"
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

func (r *UserReaderRepository) FindUserById(ctx context.Context, id string) (*entity.User, error) {
	var userModel model.User

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constant.ErrUserNotFound
		}

		return nil, err
	}

	user := entity.User{
		ID:       userModel.ID.String(),
		Email:    userModel.Email,
		Password: userModel.Password,
		Fullname: userModel.Fullname,
	}

	return &user, nil
}

func (r *UserReaderRepository) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var userModel model.User

	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constant.ErrUserNotFound
		}

		return nil, err
	}

	user := entity.User{
		ID:       userModel.ID.String(),
		Email:    userModel.Email,
		Password: userModel.Password,
		Fullname: userModel.Fullname,
	}

	return &user, nil
}

func (r *UserReaderRepository) FindUserWithActiveSessionsById(ctx context.Context, id string) (*entity.User, error) {
	var userModel model.User

	err := r.db.WithContext(ctx).
		Preload("Sessions", func(db *gorm.DB) *gorm.DB {
			return db.Where("revoked_at IS NULL")
		}).
		Where("id = ?", id).
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

	for _, sessionModel := range userModel.Sessions {
		session := &entity.Session{
			ID:           sessionModel.ID.String(),
			UserId:       sessionModel.UserId.String(),
			AccessToken:  sessionModel.AccessToken,
			RefreshToken: sessionModel.RefreshToken,
			RevokedAt:    sessionModel.RevokedAt.Ptr(),
		}

		user.AddActiveSession(session)
	}

	return &user, nil
}
