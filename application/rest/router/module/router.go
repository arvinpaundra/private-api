package module

import (
	"github.com/arvinpaundra/private-api/application/rest/handler"
	"github.com/arvinpaundra/private-api/application/rest/middleware"
	"github.com/arvinpaundra/private-api/core/validator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ModuleRouter struct {
	db     *gorm.DB
	logger *zap.Logger
	vld    *validator.Validator
}

func NewModuleRouter(
	db *gorm.DB,
	logger *zap.Logger,
	vld *validator.Validator,
) *ModuleRouter {
	return &ModuleRouter{
		db:     db,
		logger: logger,
		vld:    vld,
	}
}

func (r *ModuleRouter) Private(g *gin.RouterGroup) {
	h := handler.NewModuleHandler(r.db, r.logger, r.vld)
	m := middleware.NewAuthenticate(r.db)

	module := g.Group("/modules", m.Authenticate())
	{
		module.POST("", h.CreateModule)
		module.GET("", h.FindAllModules)
	}

	moduleDetail := module.Group("/:module_slug")
	{
		moduleDetail.GET("", h.FindDetailModule)
		moduleDetail.DELETE("", h.DeleteModule)
		moduleDetail.GET("/questions", h.FindDetailModuleQuestions)
		moduleDetail.PATCH("/publish", h.TogglePublishModule)

		question := moduleDetail.Group("/questions")

		question.POST("", h.AddQuestions)
	}
}

func (r *ModuleRouter) Public(g *gin.RouterGroup) {
	h := handler.NewModuleHandler(r.db, r.logger, r.vld)

	module := g.Group("/modules/:module_slug")
	{
		module.GET("/published", h.FindPublishedModule)

		question := module.Group("/questions")

		question.GET("/:question_slug", h.FindPublishedQuestion)
	}
}
