package handler

import (
	"net/http"

	"github.com/arvinpaundra/private-api/core/format"
	"github.com/arvinpaundra/private-api/core/validator"
	"github.com/arvinpaundra/private-api/domain/grade/constant"
	"github.com/arvinpaundra/private-api/domain/grade/service"
	"github.com/arvinpaundra/private-api/infrastructure/grade"
	"github.com/arvinpaundra/private-api/infrastructure/shared"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type GradeHandler struct {
	db     *gorm.DB
	logger *zap.Logger
	vld    *validator.Validator
}

func NewGradeHandler(
	db *gorm.DB,
	logger *zap.Logger,
	vld *validator.Validator,
) *GradeHandler {
	return &GradeHandler{
		db:     db,
		logger: logger.With(zap.String("domain", "grade")),
		vld:    vld,
	}
}

func (h *GradeHandler) CreateGrade(c *gin.Context) {
	var command service.CreateGradeCommand

	err := c.ShouldBindJSON(&command)
	if err != nil {
		c.JSON(http.StatusBadRequest, format.UnprocessableEntity(err.Error()))
	}

	verrs := h.vld.Validate(command)
	if verrs != nil {
		c.JSON(http.StatusBadRequest, format.BadRequest("invalid request body", verrs))
		return
	}

	svc := service.NewCreateGrade(
		shared.NewAuthStorage(c),
		grade.NewGradeReaderRepository(h.db),
		grade.NewGradeWriterRepository(h.db),
	)

	err = svc.Execute(c.Request.Context(), &command)
	if err != nil {
		h.logger.Error("failed to create grade", zap.Error(err))

		switch err {
		case constant.ErrGradeAlreadyExists:
			c.JSON(http.StatusConflict, format.Conflict(err.Error()))
			return
		default:
			c.JSON(http.StatusInternalServerError, format.InternalServerError())
			return
		}
	}

	c.JSON(http.StatusCreated, format.SuccessCreated("grade created successfully", nil))
}

func (h *GradeHandler) UpdateGrade(c *gin.Context) {
	var command service.UpdateGradeCommand

	err := c.ShouldBindJSON(&command)
	if err != nil {
		c.JSON(http.StatusBadRequest, format.UnprocessableEntity(err.Error()))
	}

	command.ID = c.Param("id")

	verrs := h.vld.Validate(command)
	if verrs != nil {
		c.JSON(http.StatusBadRequest, format.BadRequest("invalid request body", verrs))
		return
	}

	svc := service.NewUpdateGrade(
		shared.NewAuthStorage(c),
		grade.NewGradeReaderRepository(h.db),
		grade.NewGradeWriterRepository(h.db),
	)

	err = svc.Execute(c.Request.Context(), &command)
	if err != nil {
		h.logger.Error("failed to update grade", zap.Error(err))

		switch err {
		case constant.ErrGradeNotFound:
			c.JSON(http.StatusNotFound, format.NotFound(err.Error()))
			return
		case constant.ErrGradeAlreadyExists:
			c.JSON(http.StatusConflict, format.Conflict(err.Error()))
			return
		default:
			c.JSON(http.StatusInternalServerError, format.InternalServerError())
			return
		}
	}

	c.JSON(http.StatusOK, format.SuccessOK("grade updated successfully", nil))
}

func (h *GradeHandler) FindDetailGrade(c *gin.Context) {
	command := service.FindDetailGradeCommand{
		ID: c.Param("id"),
	}

	svc := service.NewFindDetailGrade(
		shared.NewAuthStorage(c),
		grade.NewGradeReaderRepository(h.db),
	)

	result, err := svc.Execute(c.Request.Context(), &command)
	if err != nil {
		h.logger.Error("failed to find detail grade", zap.Error(err))

		switch err {
		case constant.ErrGradeNotFound:
			c.JSON(http.StatusNotFound, format.NotFound(err.Error()))
			return
		default:
			c.JSON(http.StatusInternalServerError, format.InternalServerError())
			return
		}
	}

	c.JSON(http.StatusOK, format.SuccessOK("grade detail fetched successfully", result))
}

func (h *GradeHandler) FindAllGrades(c *gin.Context) {
	var command service.FindAllGradesCommand

	err := c.ShouldBindQuery(&command)
	if err != nil {
		c.JSON(http.StatusBadRequest, format.UnprocessableEntity(err.Error()))
		return
	}

	svc := service.NewFindAllGrades(
		shared.NewAuthStorage(c),
		grade.NewGradeReaderRepository(h.db),
	)

	result, err := svc.Execute(c.Request.Context(), &command)
	if err != nil {
		h.logger.Error("failed to find all grades", zap.Error(err))

		c.JSON(http.StatusInternalServerError, format.InternalServerError())
		return
	}

	c.JSON(http.StatusOK, format.SuccessOK("grades fetched successfully", result))
}

func (h *GradeHandler) DeleteGrade(c *gin.Context) {
	command := service.DeleteGradeCommand{
		ID: c.Param("id"),
	}

	svc := service.NewDeleteGrade(
		shared.NewAuthStorage(c),
		grade.NewGradeReaderRepository(h.db),
		grade.NewGradeWriterRepository(h.db),
	)

	err := svc.Execute(c.Request.Context(), &command)
	if err != nil {
		h.logger.Error("failed to delete grade", zap.Error(err))

		switch err {
		case constant.ErrGradeNotFound:
			c.JSON(http.StatusNotFound, format.NotFound(err.Error()))
			return
		default:
			c.JSON(http.StatusInternalServerError, format.InternalServerError())
			return
		}
	}

	c.JSON(http.StatusOK, format.SuccessOK("grade deleted successfully", nil))
}
