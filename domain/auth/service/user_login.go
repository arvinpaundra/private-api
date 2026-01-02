package service

import (
	"context"
	"errors"

	"github.com/arvinpaundra/private-api/core/token"
	"github.com/arvinpaundra/private-api/core/util"
	"github.com/arvinpaundra/private-api/domain/auth/constant"
	"github.com/arvinpaundra/private-api/domain/auth/entity"
	"github.com/arvinpaundra/private-api/domain/auth/repository"
	"github.com/arvinpaundra/private-api/domain/auth/response"
)

type UserLoginCommand struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserLogin struct {
	userReader repository.UserReader
	userWriter repository.UserWriter
	tokenable  token.Tokenable
	uow        repository.UnitOfWork
}

func NewUserLogin(
	userReader repository.UserReader,
	userWriter repository.UserWriter,
	tokenable token.Tokenable,
	uow repository.UnitOfWork,
) *UserLogin {
	return &UserLogin{
		userReader: userReader,
		userWriter: userWriter,
		tokenable:  tokenable,
		uow:        uow,
	}
}

func (s *UserLogin) Execute(ctx context.Context, command UserLoginCommand) (*response.UserLogin, error) {
	user, err := s.userReader.FindUserByEmail(ctx, command.Email)
	if err != nil {
		if errors.Is(err, constant.ErrUserNotFound) {
			return nil, constant.ErrWrongEmailOrPassword
		}
		return nil, err
	}

	err = util.CompareHashAndString(user.Password, command.Password)
	if err != nil {
		return nil, constant.ErrWrongEmailOrPassword
	}

	accessToken, err := s.tokenable.Encode(user.ID, constant.TokenValidSevenDays, constant.TokenValidImmediately)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenable.Encode(user.ID, constant.TokenValidSevenDays, constant.TokenValidImmediately)
	if err != nil {
		return nil, err
	}

	session := entity.NewSession(user.ID, accessToken, refreshToken)

	user.AddActiveSession(session)

	user.MarkUpdate()

	tx, err := s.uow.Begin()
	if err != nil {
		return nil, err
	}

	err = tx.UserWriter().Save(ctx, user)
	if err != nil {
		if uowErr := tx.Rollback(); uowErr != nil {
			return nil, uowErr
		}

		return nil, err
	}

	uowErr := tx.Commit()
	if uowErr != nil {
		return nil, uowErr
	}

	result := &response.UserLogin{
		UserId:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return result, nil
}
