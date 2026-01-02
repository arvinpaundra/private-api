package service

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/user/repository"
	"github.com/arvinpaundra/private-api/domain/user/response"
)

type FindUserByID struct {
	userReader repository.UserReader
}

func NewFindUserByID(userReader repository.UserReader) *FindUserByID {
	return &FindUserByID{
		userReader: userReader,
	}
}

func (s *FindUserByID) Execute(ctx context.Context, userID string) (*response.User, error) {
	user, err := s.userReader.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &response.User{
		ID:       user.ID,
		Email:    user.Email,
		Fullname: user.Fullname,
	}, nil
}
