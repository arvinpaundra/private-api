package auth

import (
	"context"

	"github.com/arvinpaundra/private-api/core/util"
	"github.com/arvinpaundra/private-api/domain/auth/entity"
	"github.com/arvinpaundra/private-api/domain/auth/repository"
	"github.com/arvinpaundra/private-api/model"
	"github.com/guregu/null/v6"
	"gorm.io/gorm"
)

var _ repository.UserWriter = (*UserWriterRepository)(nil)

type UserWriterRepository struct {
	db *gorm.DB
}

func NewUserWriterRepository(db *gorm.DB) *UserWriterRepository {
	return &UserWriterRepository{db: db}
}

func (r *UserWriterRepository) Save(ctx context.Context, user *entity.User) error {
	if user.IsCreated() {
		return r.insert(ctx, user)
	} else if user.IsUpdated() {
		return r.update(ctx, user)
	}

	return nil
}

func (r *UserWriterRepository) insert(ctx context.Context, user *entity.User) error {
	userModel := model.User{
		ID:       util.ParseUUID(user.ID),
		Email:    user.Email,
		Password: user.Password,
		Fullname: user.Fullname,
	}

	if err := r.db.WithContext(ctx).Create(&userModel).Error; err != nil {
		return err
	}

	return nil
}

func (r *UserWriterRepository) update(ctx context.Context, user *entity.User) error {
	userModel := model.User{
		Email:    user.Email,
		Password: user.Password,
		Fullname: user.Fullname,
	}

	if err := r.db.WithContext(ctx).Where("id = ?", user.ID).Updates(&userModel).Error; err != nil {
		return err
	}

	for _, session := range user.ActiveSessions {
		sessionModel := model.Session{
			ID:           util.ParseUUID(session.ID),
			UserId:       util.ParseUUID(user.ID),
			AccessToken:  session.AccessToken,
			RefreshToken: session.RefreshToken,
			RevokedAt:    null.TimeFromPtr(session.RevokedAt),
		}

		if session.IsCreated() {
			if err := r.db.WithContext(ctx).Create(&sessionModel).Error; err != nil {
				return err
			}
		} else if session.IsUpdated() {
			if err := r.db.WithContext(ctx).Where("id = ?", session.ID).Updates(&sessionModel).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
