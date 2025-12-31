package submission

import (
	"context"
	"strings"

	"github.com/arvinpaundra/private-api/domain/module/service"
	"github.com/arvinpaundra/private-api/domain/submission/constant"
	"github.com/arvinpaundra/private-api/domain/submission/entity"
	"github.com/arvinpaundra/private-api/domain/submission/repository"
	"github.com/arvinpaundra/private-api/infrastructure/module"
	"gorm.io/gorm"
)

var _ repository.ModuleACL = (*ModuleACLAdapter)(nil)

type ModuleACLAdapter struct {
	db *gorm.DB
}

func NewModuleACLAdapter(db *gorm.DB) *ModuleACLAdapter {
	return &ModuleACLAdapter{
		db: db,
	}
}

func (a *ModuleACLAdapter) GetCorrectAnswer(ctx context.Context, moduleSlug, questionSlug string) (*entity.Choice, error) {
	svc := service.NewGetCorrectAnswer(
		module.NewModuleReaderRepository(a.db),
	)

	correctChoice, err := svc.Execute(ctx, &service.GetCorrectAnswerCommand{
		ModuleSlug:   moduleSlug,
		QuestionSlug: questionSlug,
	})
	if err != nil {
		return nil, err
	}

	// Map to submission domain entity
	return &entity.Choice{
		ID:              correctChoice.ID,
		Content:         correctChoice.Content,
		IsCorrectAnswer: correctChoice.IsCorrectAnswer,
	}, nil
}

func (a *ModuleACLAdapter) GetNextQuestionSlug(ctx context.Context, moduleSlug, currentQuestionSlug string) (*string, error) {
	svc := service.NewFindPublishedQuestion(
		module.NewModuleReaderRepository(a.db),
	)

	questionDetail, err := svc.Execute(ctx, &service.FindPublishedQuestionCommand{
		ModuleSlug:   moduleSlug,
		QuestionSlug: currentQuestionSlug,
	})
	if err != nil {
		return nil, err
	}

	// Return next question slug
	return questionDetail.NextQuestionSlug, nil
}

func (a *ModuleACLAdapter) GetPublishedModule(ctx context.Context, moduleSlug string) (*entity.Module, error) {
	svc := service.NewValidatePublishedModule(
		module.NewModuleReaderRepository(a.db),
	)

	module, err := svc.Execute(ctx, &service.ValidatePublishedModuleCommand{
		ModuleSlug: moduleSlug,
	})
	if err != nil {
		if strings.Contains(err.Error(), constant.ErrModuleNotFound.Error()) {
			return nil, constant.ErrModuleNotFound
		}
		return nil, err
	}

	return &entity.Module{
		ID:   module.ID,
		Slug: module.Slug,
	}, nil
}

func (a *ModuleACLAdapter) GetQuestionBySlug(ctx context.Context, moduleSlug, questionSlug string) (*entity.Question, error) {
	// Use module domain service to get question details
	moduleService := service.NewFindPublishedQuestion(
		module.NewModuleReaderRepository(a.db),
	)

	questionDetail, err := moduleService.Execute(ctx, &service.FindPublishedQuestionCommand{
		ModuleSlug:   moduleSlug,
		QuestionSlug: questionSlug,
	})
	if err != nil {
		if strings.Contains(err.Error(), constant.ErrModuleNotFound.Error()) {
			return nil, constant.ErrModuleNotFound
		} else if strings.Contains(err.Error(), constant.ErrQuestionNotFound.Error()) {
			return nil, constant.ErrQuestionNotFound
		}
		return nil, err
	}

	// Map choices to submission domain
	choices := make([]*entity.Choice, len(questionDetail.Choices))
	for i, choice := range questionDetail.Choices {
		choices[i] = &entity.Choice{
			ID:      choice.ID,
			Content: choice.Content,
		}
	}

	// Map to submission domain entity
	return &entity.Question{
		ID:               questionDetail.ID,
		Content:          questionDetail.Content,
		Slug:             questionDetail.Slug,
		NextQuestionSlug: questionDetail.NextQuestionSlug,
		Choices:          choices,
	}, nil
}

func (a *ModuleACLAdapter) GetTotalQuestions(ctx context.Context, moduleSlug string) (int, error) {
	svc := service.NewCountModuleQuestions(
		module.NewModuleReaderRepository(a.db),
	)

	total, err := svc.Execute(ctx, &service.CountModuleQuestionsCommand{
		ModuleSlug: moduleSlug,
	})
	if err != nil {
		return 0, err
	}

	return total, nil
}
