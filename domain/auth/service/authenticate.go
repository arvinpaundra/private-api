package service

import (
	"context"

	"github.com/arvinpaundra/private-api/core/token"
	"github.com/arvinpaundra/private-api/domain/auth/repository"
	"github.com/arvinpaundra/private-api/domain/auth/response"
)

type CheckSessionCommand struct {
	AccessToken string `json:"access_token" binding:"required"`
}

type CheckSession struct {
	userReader repository.UserReader
	tokenable  token.Tokenable
}

func NewCheckSession(
	userReader repository.UserReader,
	tokenable token.Tokenable,
) CheckSession {
	return CheckSession{
		userReader: userReader,
		tokenable:  tokenable,
	}
}

func (s *CheckSession) Execute(ctx context.Context, command CheckSessionCommand) (*response.AuthenticatedUser, error) {
	claims, err := s.tokenable.Decode(command.AccessToken)
	if err != nil {
		return nil, err
	}

	user, err := s.userReader.FindUserById(ctx, claims.Identifier)
	if err != nil {
		return nil, err
	}

	result := &response.AuthenticatedUser{
		UserId:   user.ID,
		Email:    user.Email,
		Fullname: user.Fullname,
	}

	return result, nil
}
