package handler

import (
	"net/http"

	"github.com/arvinpaundra/private-api/core/format"
	"github.com/arvinpaundra/private-api/core/validator"
	"github.com/arvinpaundra/private-api/domain/submission/constant"
	"github.com/arvinpaundra/private-api/domain/submission/service"
	"github.com/arvinpaundra/private-api/infrastructure/submission"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SubmissionHandler struct {
	db     *gorm.DB
	logger *zap.Logger
	vld    *validator.Validator
}

func NewSubmissionHandler(
	db *gorm.DB,
	logger *zap.Logger,
	vld *validator.Validator,
) *SubmissionHandler {
	return &SubmissionHandler{
		db:     db,
		logger: logger.With(zap.String("domain", "submission")),
		vld:    vld,
	}
}

func (h *SubmissionHandler) StartSubmission(c *gin.Context) {
	var command service.StartSubmissionCommand

	err := c.ShouldBindJSON(&command)
	if err != nil {
		c.JSON(http.StatusBadRequest, format.UnprocessableEntity(err.Error()))
		return
	}

	command.ModuleSlug = c.Param("module_slug")

	verrs := h.vld.Validate(command)
	if verrs != nil {
		c.JSON(http.StatusBadRequest, format.BadRequest("invalid request body", verrs))
		return
	}

	svc := service.NewStartSubmission(
		submission.NewUnitOfWork(h.db),
		submission.NewModuleACLAdapter(h.db),
	)

	result, err := svc.Execute(c.Request.Context(), &command)
	if err != nil {
		h.logger.Error("failed to start submission", zap.Error(err))

		switch err {
		case constant.ErrModuleNotFound:
			c.JSON(http.StatusNotFound, format.NotFound(err.Error()))
			return
		default:
			c.JSON(http.StatusInternalServerError, format.InternalServerError())
			return
		}
	}

	c.JSON(http.StatusCreated, format.SuccessCreated("submission started successfully", result))
}

func (h *SubmissionHandler) SubmitAnswer(c *gin.Context) {
	var command service.SubmitAnswerCommand

	err := c.ShouldBindJSON(&command)
	if err != nil {
		c.JSON(http.StatusBadRequest, format.UnprocessableEntity(err.Error()))
		return
	}

	command.ModuleSlug = c.Param("module_slug")
	command.SubmissionCode = c.Param("submission_code")

	verrs := h.vld.Validate(command)
	if verrs != nil {
		c.JSON(http.StatusBadRequest, format.BadRequest("invalid request body", verrs))
		return
	}

	svc := service.NewSubmitAnswer(
		submission.NewSubmissionReaderRepository(h.db),
		submission.NewUnitOfWork(h.db),
		submission.NewModuleACLAdapter(h.db),
	)

	result, err := svc.Execute(c.Request.Context(), &command)
	if err != nil {
		h.logger.Error("failed to submit answer", zap.Error(err))

		switch err {
		case constant.ErrSubmissionNotFound, constant.ErrModuleNotFound, constant.ErrQuestionNotFound, constant.ErrChoiceNotFound:
			c.JSON(http.StatusNotFound, format.NotFound(err.Error()))
			return
		case constant.ErrSubmissionAlreadyDone:
			c.JSON(http.StatusBadRequest, format.BadRequest(err.Error(), nil))
			return
		case constant.ErrDuplicateAnswer:
			c.JSON(http.StatusConflict, format.BadRequest(err.Error(), nil))
			return
		default:
			c.JSON(http.StatusInternalServerError, format.InternalServerError())
			return
		}
	}

	c.JSON(http.StatusOK, format.SuccessOK("answer submitted successfully", result))
}

func (h *SubmissionHandler) FinalizeSubmission(c *gin.Context) {
	command := service.FinalizeSubmissionCommand{
		ModuleSlug:     c.Param("module_slug"),
		SubmissionCode: c.Param("submission_code"),
	}

	verrs := h.vld.Validate(command)
	if verrs != nil {
		c.JSON(http.StatusBadRequest, format.BadRequest("invalid request", verrs))
		return
	}

	svc := service.NewFinalizeSubmission(
		submission.NewSubmissionReaderRepository(h.db),
		submission.NewModuleACLAdapter(h.db),
		submission.NewUnitOfWork(h.db),
	)

	result, err := svc.Execute(c.Request.Context(), &command)
	if err != nil {
		h.logger.Error("failed to finalize submission", zap.Error(err))

		switch err {
		case constant.ErrSubmissionNotFound, constant.ErrModuleNotFound:
			c.JSON(http.StatusNotFound, format.NotFound(err.Error()))
			return
		case constant.ErrSubmissionAlreadyDone:
			c.JSON(http.StatusBadRequest, format.BadRequest(err.Error(), nil))
			return
		default:
			c.JSON(http.StatusInternalServerError, format.InternalServerError())
			return
		}
	}

	c.JSON(http.StatusOK, format.SuccessOK("submission finalized successfully", result))
}

func (h *SubmissionHandler) GetAllSubmissions(c *gin.Context) {
	var query service.FindAllSubmissionQuery

	err := c.ShouldBindQuery(&query)
	if err != nil {
		c.JSON(http.StatusBadRequest, format.UnprocessableEntity(err.Error()))
		return
	}

	svc := service.NewFindAllSubmission(
		submission.NewSubmissionReaderRepository(h.db),
		submission.NewModuleACLAdapter(h.db),
	)

	result, err := svc.Execute(c.Request.Context(), &query)
	if err != nil {
		h.logger.Error("failed to get all submissions", zap.Error(err))

		c.JSON(http.StatusInternalServerError, format.InternalServerError())
		return
	}

	c.JSON(http.StatusOK, format.SuccessOK("submissions retrieved successfully", result))
}
