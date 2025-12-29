package module

import (
	"github.com/arvinpaundra/private-api/application/rest/handler"
	"github.com/arvinpaundra/private-api/application/rest/middleware"
	"github.com/arvinpaundra/private-api/core/validator"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ModuleRouter struct {
	db  *gorm.DB
	vld *validator.Validator
}

func NewModuleRouter(db *gorm.DB, vld *validator.Validator) *ModuleRouter {
	return &ModuleRouter{
		db:  db,
		vld: vld,
	}
}

func (r *ModuleRouter) Private(g *gin.RouterGroup) {
	h := handler.NewModuleHandler(r.db, r.vld)
	m := middleware.NewAuthenticate(r.db)

	module := g.Group("/modules", m.Authenticate())
	{
		module.POST("", h.CreateModule)
		module.GET("", h.FindAllModules)
	}

	moduleDetail := module.Group("/:module_slug")
	{
		moduleDetail.GET("", h.FindDetailModule)
		moduleDetail.PATCH("/publish", h.TogglePublishModule)

		question := moduleDetail.Group("/questions")

		question.POST("", h.AddQuestions)
	}
}
