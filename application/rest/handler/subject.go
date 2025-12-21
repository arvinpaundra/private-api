package handler

import (
	"net/http"

	"github.com/arvinpaundra/private-api/core/format"
	"github.com/arvinpaundra/private-api/core/validator"
	"github.com/arvinpaundra/private-api/domain/subject/constant"
	"github.com/arvinpaundra/private-api/domain/subject/service"
	"github.com/arvinpaundra/private-api/infrastructure/shared"
	"github.com/arvinpaundra/private-api/infrastructure/subject"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SubjectHandler struct {
	db  *gorm.DB
	vld *validator.Validator
}

func NewSubjectHandler(db *gorm.DB, vld *validator.Validator) *SubjectHandler {
	return &SubjectHandler{
		db:  db,
		vld: vld,
	}
}

func (h *SubjectHandler) CreateSubject(c *gin.Context) {
	var command service.CreateSubjectCommand

	err := c.ShouldBindJSON(&command)
	if err != nil {
		c.JSON(http.StatusBadRequest, format.UnprocessableEntity(err.Error()))
	}

	verrs := h.vld.Validate(command)
	if verrs != nil {
		c.JSON(http.StatusBadRequest, format.BadRequest("invalid request body", verrs))
		return
	}

	svc := service.NewCreateSubject(
		shared.NewAuthStorage(c),
		subject.NewSubjectReaderRepository(h.db),
		subject.NewSubjectWriterRepository(h.db),
	)

	err = svc.Execute(c.Request.Context(), &command)
	if err != nil {
		switch err {
		case constant.ErrSubjectAlreadyExists:
			c.JSON(http.StatusConflict, format.Conflict(err.Error()))
			return
		default:
			c.JSON(http.StatusInternalServerError, format.InternalServerError())
			return
		}
	}

	c.JSON(http.StatusCreated, format.SuccessCreated("subject created successfully", nil))
}

func (h *SubjectHandler) UpdateSubject(c *gin.Context) {
	var command service.UpdateSubjectCommand

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

	svc := service.NewUpdateSubject(
		shared.NewAuthStorage(c),
		subject.NewSubjectReaderRepository(h.db),
		subject.NewSubjectWriterRepository(h.db),
	)

	err = svc.Execute(c.Request.Context(), &command)
	if err != nil {
		switch err {
		case constant.ErrSubjectNotFound:
			c.JSON(http.StatusNotFound, format.NotFound(err.Error()))
			return
		case constant.ErrSubjectAlreadyExists:
			c.JSON(http.StatusConflict, format.Conflict(err.Error()))
			return
		default:
			c.JSON(http.StatusInternalServerError, format.InternalServerError())
			return
		}
	}

	c.JSON(http.StatusOK, format.SuccessOK("subject updated successfully", nil))
}

func (h *SubjectHandler) FindDetailSubject(c *gin.Context) {
	command := service.FindDetailSubjectCommand{
		ID: c.Param("id"),
	}

	svc := service.NewFindDetailSubject(
		shared.NewAuthStorage(c),
		subject.NewSubjectReaderRepository(h.db),
	)

	result, err := svc.Execute(c.Request.Context(), &command)
	if err != nil {
		switch err {
		case constant.ErrSubjectNotFound:
			c.JSON(http.StatusNotFound, format.NotFound(err.Error()))
			return
		default:
			c.JSON(http.StatusInternalServerError, format.InternalServerError())
			return
		}
	}

	c.JSON(http.StatusOK, format.SuccessOK("subject detail fetched successfully", result))
}

func (h *SubjectHandler) FindAllSubjects(c *gin.Context) {
	var command service.FindAllSubjectsCommand

	err := c.ShouldBindQuery(&command)
	if err != nil {
		c.JSON(http.StatusBadRequest, format.UnprocessableEntity(err.Error()))
		return
	}

	svc := service.NewFindAllSubjects(
		shared.NewAuthStorage(c),
		subject.NewSubjectReaderRepository(h.db),
	)

	result, err := svc.Execute(c.Request.Context(), &command)
	if err != nil {
		c.JSON(http.StatusInternalServerError, format.InternalServerError())
		return
	}

	c.JSON(http.StatusOK, format.SuccessOK("subjects fetched successfully", result))
}

func (h *SubjectHandler) DeleteSubject(c *gin.Context) {
	command := service.DeleteSubjectCommand{
		ID: c.Param("id"),
	}

	svc := service.NewDeleteSubject(
		shared.NewAuthStorage(c),
		subject.NewSubjectReaderRepository(h.db),
		subject.NewSubjectWriterRepository(h.db),
	)

	err := svc.Execute(c.Request.Context(), &command)
	if err != nil {
		switch err {
		case constant.ErrSubjectNotFound:
			c.JSON(http.StatusNotFound, format.NotFound(err.Error()))
			return
		default:
			c.JSON(http.StatusInternalServerError, format.InternalServerError())
			return
		}
	}

	c.JSON(http.StatusOK, format.SuccessOK("subject deleted successfully", nil))
}
