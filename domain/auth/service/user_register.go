package service

import (
	"context"
	"errors"

	"github.com/arvinpaundra/private-api/core/util"
	"github.com/arvinpaundra/private-api/domain/auth/constant"
	"github.com/arvinpaundra/private-api/domain/auth/entity"
	"github.com/arvinpaundra/private-api/domain/auth/repository"
)

type UserRegisterCommand struct {
	Email    string `json:"email" validate:"required,email,max=50"`
	Password string `json:"password" validate:"required,min=8"`
	Fullname string `json:"fullname" validate:"required,min=3,max=100"`
}

type UserRegister struct {
	userReader repository.UserReader
	userWriter repository.UserWriter
	uow        repository.UnitOfWork
}

func NewUserRegister(
	userReader repository.UserReader,
	userWriter repository.UserWriter,
	uow repository.UnitOfWork,
) *UserRegister {
	return &UserRegister{
		userReader: userReader,
		userWriter: userWriter,
		uow:        uow,
	}
}

func (s *UserRegister) Execute(ctx context.Context, command UserRegisterCommand) error {
	existingUser, err := s.userReader.FindUserByEmail(ctx, command.Email)
	if err != nil && !errors.Is(err, constant.ErrUserNotFound) {
		return err
	}

	if !existingUser.IsEmpty() {
		return constant.ErrEmailAlreadyExists
	}

	password, err := util.HashString(command.Password)
	if err != nil {
		return err
	}

	user := entity.NewUser(
		command.Email,
		password,
		command.Fullname,
	)

	tx, err := s.uow.Begin()
	if err != nil {
		return err
	}

	if err := tx.UserWriter().Save(ctx, user); err != nil {
		if uowErr := tx.Rollback(); uowErr != nil {
			return uowErr
		}

		return err
	}

	uowErr := tx.Commit()
	if uowErr != nil {
		return uowErr
	}

	return nil
}
