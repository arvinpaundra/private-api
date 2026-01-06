package handler

import (
	"net/http"

	"github.com/arvinpaundra/private-api/core/format"
	"github.com/arvinpaundra/private-api/core/validator"
	"github.com/arvinpaundra/private-api/domain/module/constant"
	"github.com/arvinpaundra/private-api/domain/module/service"
	"github.com/arvinpaundra/private-api/infrastructure/module"
	"github.com/arvinpaundra/private-api/infrastructure/shared"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ModuleHandler struct {
	db     *gorm.DB
	logger *zap.Logger
	vld    *validator.Validator
}

func NewModuleHandler(
	db *gorm.DB,
	logger *zap.Logger,
	vld *validator.Validator,
) *ModuleHandler {
	return &ModuleHandler{
		db:     db,
		logger: logger.With(zap.String("domain", "module")),
		vld:    vld,
	}
}

func (h *ModuleHandler) CreateModule(c *gin.Context) {
	var command service.CreateModuleCommand

	err := c.ShouldBindJSON(&command)
	if err != nil {
		c.JSON(http.StatusBadRequest, format.UnprocessableEntity(err.Error()))
		return
	}

	verrs := h.vld.Validate(command)
	if verrs != nil {
		c.JSON(http.StatusBadRequest, format.BadRequest("invalid request body", verrs))
		return
	}

	svc := service.NewCreateModule(
		shared.NewAuthStorage(c),
		module.NewModuleWriterRepository(h.db),
		module.NewSubjectACLAdapter(h.db, shared.NewAuthStorage(c)),
		module.NewGradeACLAdapter(h.db, shared.NewAuthStorage(c)),
	)

	slug, err := svc.Execute(c.Request.Context(), &command)
	if err != nil {
		h.logger.Error("failed to create module", zap.Error(err))

		switch err {
		case constant.ErrSubjectNotFound, constant.ErrGradeNotFound, constant.ErrModuleNotFound:
			c.JSON(http.StatusNotFound, format.NotFound(err.Error()))
			return
		default:
			c.JSON(http.StatusInternalServerError, format.InternalServerError())
			return
		}
	}

	c.JSON(http.StatusCreated, format.SuccessCreated("module created successfully", gin.H{
		"slug": slug,
	}))
}

func (h *ModuleHandler) FindAllModules(c *gin.Context) {
	var command service.FindAllModulesCommand

	err := c.ShouldBindQuery(&command)
	if err != nil {
		c.JSON(http.StatusBadRequest, format.UnprocessableEntity(err.Error()))
		return
	}

	svc := service.NewFindAllModules(
		shared.NewAuthStorage(c),
		module.NewModuleReaderRepository(h.db),
		module.NewSubjectACLAdapter(h.db, shared.NewAuthStorage(c)),
		module.NewGradeACLAdapter(h.db, shared.NewAuthStorage(c)),
	)

	result, err := svc.Execute(c.Request.Context(), &command)
	if err != nil {
		h.logger.Error("failed to find all modules", zap.Error(err))

		switch err {
		case constant.ErrSubjectNotFound, constant.ErrGradeNotFound, constant.ErrModuleNotFound:
			c.JSON(http.StatusNotFound, format.NotFound(err.Error()))
			return
		default:
			c.JSON(http.StatusInternalServerError, format.InternalServerError())
			return
		}
	}

	c.JSON(http.StatusOK, format.SuccessOK("modules retrieved successfully", result))
}

func (h *ModuleHandler) AddQuestions(c *gin.Context) {
	moduleSlug := c.Param("module_slug")

	var command service.AddQuestionsCommand

	err := c.ShouldBindJSON(&command)
	if err != nil {
		c.JSON(http.StatusBadRequest, format.UnprocessableEntity(err.Error()))
		return
	}

	// Set module slug from URL param
	command.ModuleSlug = moduleSlug

	verrs := h.vld.Validate(command)
	if verrs != nil {
		c.JSON(http.StatusBadRequest, format.BadRequest("invalid request body", verrs))
		return
	}

	svc := service.NewAddQuestions(
		shared.NewAuthStorage(c),
		module.NewModuleReaderRepository(h.db),
		module.NewUnitOfWork(h.db),
	)

	err = svc.Execute(c.Request.Context(), &command)
	if err != nil {
		h.logger.Error("failed to add questions", zap.Error(err))

		switch err {
		case constant.ErrModuleNotFound, constant.ErrQuestionNotFound:
			c.JSON(http.StatusNotFound, format.NotFound(err.Error()))
			return
		case constant.ErrMinTwoChoices, constant.ErrMaxFourChoices, constant.ErrMultipleCorrectAnswers, constant.ErrNoCorrectAnswer:
			c.JSON(http.StatusBadRequest, format.BadRequest(err.Error(), nil))
			return
		default:
			c.JSON(http.StatusInternalServerError, format.InternalServerError())
			return
		}
	}

	c.JSON(http.StatusCreated, format.SuccessCreated("question(s) added successfully", nil))
}

func (h *ModuleHandler) FindDetailModule(c *gin.Context) {
	slug := c.Param("module_slug")

	command := service.FindDetailModuleCommand{
		Slug: slug,
	}

	svc := service.NewFindDetailModule(
		shared.NewAuthStorage(c),
		module.NewModuleReaderRepository(h.db),
	)

	result, err := svc.Execute(c.Request.Context(), &command)
	if err != nil {
		h.logger.Error("failed to find detail module", zap.Error(err))

		switch err {
		case constant.ErrModuleNotFound:
			c.JSON(http.StatusNotFound, format.NotFound(err.Error()))
			return
		default:
			c.JSON(http.StatusInternalServerError, format.InternalServerError())
			return
		}
	}

	c.JSON(http.StatusOK, format.SuccessOK("module retrieved successfully", result))
}

func (h *ModuleHandler) FindPublishedModule(c *gin.Context) {
	slug := c.Param("module_slug")

	command := service.FindPublishedModuleCommand{
		Slug: slug,
	}

	svc := service.NewFindPublishedModule(
		module.NewModuleReaderRepository(h.db),
	)

	result, err := svc.Execute(c.Request.Context(), &command)
	if err != nil {
		h.logger.Error("failed to find published module", zap.Error(err))

		switch err {
		case constant.ErrModuleNotFound:
			c.JSON(http.StatusNotFound, format.NotFound(err.Error()))
			return
		default:
			c.JSON(http.StatusInternalServerError, format.InternalServerError())
			return
		}
	}

	c.JSON(http.StatusOK, format.SuccessOK("module retrieved successfully", result))
}

func (h *ModuleHandler) FindDetailModuleQuestions(c *gin.Context) {
	slug := c.Param("module_slug")

	command := service.FindDetailModuleQuestionsCommand{
		Slug: slug,
	}

	svc := service.NewFindDetailModuleQuestions(
		shared.NewAuthStorage(c),
		module.NewModuleReaderRepository(h.db),
		module.NewSubjectACLAdapter(h.db, shared.NewAuthStorage(c)),
		module.NewGradeACLAdapter(h.db, shared.NewAuthStorage(c)),
	)

	result, err := svc.Execute(c.Request.Context(), &command)
	if err != nil {
		h.logger.Error("failed to find detail module questions", zap.Error(err))

		switch err {
		case constant.ErrModuleNotFound:
			c.JSON(http.StatusNotFound, format.NotFound(err.Error()))
			return
		default:
			c.JSON(http.StatusInternalServerError, format.InternalServerError())
			return
		}
	}

	c.JSON(http.StatusOK, format.SuccessOK("module with questions retrieved successfully", result))
}

func (h *ModuleHandler) TogglePublishModule(c *gin.Context) {
	slug := c.Param("module_slug")

	command := service.TogglePublishModuleCommand{
		Slug: slug,
	}

	verrs := h.vld.Validate(command)
	if verrs != nil {
		c.JSON(http.StatusBadRequest, format.BadRequest("invalid request", verrs))
		return
	}

	svc := service.NewTogglePublishModule(
		shared.NewAuthStorage(c),
		module.NewModuleReaderRepository(h.db),
		module.NewModuleWriterRepository(h.db),
	)

	err := svc.Execute(c.Request.Context(), &command)
	if err != nil {
		h.logger.Error("failed to toggle publish module", zap.Error(err))

		switch err {
		case constant.ErrModuleNotFound:
			c.JSON(http.StatusNotFound, format.NotFound(err.Error()))
			return
		default:
			c.JSON(http.StatusInternalServerError, format.InternalServerError())
			return
		}
	}

	c.JSON(http.StatusOK, format.SuccessOK("module publish status toggled successfully", nil))
}

func (h *ModuleHandler) FindPublishedQuestion(c *gin.Context) {
	moduleSlug := c.Param("module_slug")
	questionSlug := c.Param("question_slug")

	command := service.FindPublishedQuestionCommand{
		ModuleSlug:   moduleSlug,
		QuestionSlug: questionSlug,
	}

	verrs := h.vld.Validate(command)
	if verrs != nil {
		c.JSON(http.StatusBadRequest, format.BadRequest("invalid request", verrs))
		return
	}

	svc := service.NewFindPublishedQuestion(
		module.NewModuleReaderRepository(h.db),
	)

	result, err := svc.Execute(c.Request.Context(), &command)
	if err != nil {
		h.logger.Error("failed to find published question", zap.Error(err))

		switch err {
		case constant.ErrModuleNotFound, constant.ErrQuestionNotFound:
			c.JSON(http.StatusNotFound, format.NotFound(err.Error()))
			return
		default:
			c.JSON(http.StatusInternalServerError, format.InternalServerError())
			return
		}
	}

	c.JSON(http.StatusOK, format.SuccessOK("question retrieved successfully", result))
}

func (h *ModuleHandler) DeleteModule(c *gin.Context) {
	slug := c.Param("module_slug")

	command := service.DeleteModuleCommand{
		Slug: slug,
	}

	verrs := h.vld.Validate(command)
	if verrs != nil {
		c.JSON(http.StatusBadRequest, format.BadRequest("invalid request", verrs))
		return
	}

	svc := service.NewDeleteModule(
		shared.NewAuthStorage(c),
		module.NewModuleReaderRepository(h.db),
		module.NewModuleWriterRepository(h.db),
	)

	err := svc.Execute(c.Request.Context(), &command)
	if err != nil {
		h.logger.Error("failed to delete module", zap.Error(err))

		switch err {
		case constant.ErrModuleNotFound:
			c.JSON(http.StatusNotFound, format.NotFound(err.Error()))
			return
		default:
			c.JSON(http.StatusInternalServerError, format.InternalServerError())
			return
		}
	}

	c.JSON(http.StatusOK, format.SuccessOK("module deleted successfully", nil))
}
