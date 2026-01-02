package dashboard

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/dashboard/entity"
	"github.com/arvinpaundra/private-api/domain/user/service"
	"github.com/arvinpaundra/private-api/infrastructure/user"
	"gorm.io/gorm"
)

type UserACLAdapter struct {
	db *gorm.DB
}

func NewUserACLAdapter(db *gorm.DB) *UserACLAdapter {
	return &UserACLAdapter{
		db: db,
	}
}

func (a *UserACLAdapter) GetUserByID(ctx context.Context, userID string) (*entity.User, error) {
	svc := service.NewFindUserByID(
		user.NewUserReaderRepository(a.db),
	)

	user, err := svc.Execute(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:       user.ID,
		Fullname: user.Fullname,
	}, nil
}
