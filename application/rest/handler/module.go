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
	"gorm.io/gorm"
)

type ModuleHandler struct {
	db  *gorm.DB
	vld *validator.Validator
}

func NewModuleHandler(db *gorm.DB, vld *validator.Validator) *ModuleHandler {
	return &ModuleHandler{
		db:  db,
		vld: vld,
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

	err = svc.Execute(c.Request.Context(), &command)
	if err != nil {
		switch err {
		case constant.ErrSubjectNotFound, constant.ErrGradeNotFound, constant.ErrModuleNotFound:
			c.JSON(http.StatusNotFound, format.NotFound(err.Error()))
			return
		default:
			c.JSON(http.StatusInternalServerError, format.InternalServerError())
			return
		}
	}

	c.JSON(http.StatusCreated, format.SuccessCreated("module created successfully", nil))
}
