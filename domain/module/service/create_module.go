package service

import (
	"context"

	"github.com/arvinpaundra/private-api/domain/module/constant"
	"github.com/arvinpaundra/private-api/domain/module/entity"
	"github.com/arvinpaundra/private-api/domain/module/repository"
	"github.com/arvinpaundra/private-api/domain/shared/interfaces"
)

type CreateModuleCommand struct {
	Title       string  `json:"title" validate:"required,max=100"`
	SubjectID   string  `json:"subject_id" validate:"required"`
	GradeID     string  `json:"grade_id" validate:"required"`
	Description *string `json:"description,omitempty"`
}

type CreateModule struct {
	authStorage  interfaces.AuthenticatedUser
	moduleWriter repository.ModuleWriter
	subjectACL   repository.SubjectACL
	gradeACL     repository.GradeACL
}

func NewCreateModule(
	authStorage interfaces.AuthenticatedUser,
	moduleWriter repository.ModuleWriter,
	subjectACL repository.SubjectACL,
	gradeACL repository.GradeACL,
) *CreateModule {
	return &CreateModule{
		authStorage:  authStorage,
		moduleWriter: moduleWriter,
		subjectACL:   subjectACL,
		gradeACL:     gradeACL,
	}
}

func (s *CreateModule) Execute(ctx context.Context, command *CreateModuleCommand) error {
	// check if subject exists via ACL
	isSubjectExist, err := s.subjectACL.IsSubjectExist(ctx, command.SubjectID, s.authStorage.GetUserId())
	if err != nil {
		return err
	}

	if !isSubjectExist {
		return constant.ErrSubjectNotFound
	}

	// check if grade exists via ACL
	isGradeExist, err := s.gradeACL.IsGradeExist(ctx, command.GradeID, s.authStorage.GetUserId())
	if err != nil {
		return err
	}

	if !isGradeExist {
		return constant.ErrGradeNotFound
	}

	// create module
	module := entity.NewModule(s.authStorage.GetUserId(), command.SubjectID, command.GradeID, command.Title, command.Description)

	module.GenSlug()

	// store module to persistent storage
	err = s.moduleWriter.Save(ctx, module)
	if err != nil {
		return err
	}

	return nil
}
