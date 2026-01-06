package subject

import (
	"github.com/arvinpaundra/private-api/application/rest/handler"
	"github.com/arvinpaundra/private-api/application/rest/middleware"
	"github.com/arvinpaundra/private-api/core/validator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SubjectRouter struct {
	db     *gorm.DB
	logger *zap.Logger
	vld    *validator.Validator
}

func NewSubjectRouter(
	db *gorm.DB,
	logger *zap.Logger,
	vld *validator.Validator,
) *SubjectRouter {
	return &SubjectRouter{
		db:     db,
		logger: logger,
		vld:    vld,
	}
}

func (r *SubjectRouter) Private(g *gin.RouterGroup) {
	h := handler.NewSubjectHandler(r.db, r.logger, r.vld)
	m := middleware.NewAuthenticate(r.db)

	subject := g.Group("/subjects", m.Authenticate())
	{
		subject.POST("", h.CreateSubject)
		subject.PUT("/:id", h.UpdateSubject)
		subject.DELETE("/:id", h.DeleteSubject)
		subject.GET("", h.FindAllSubjects)
		subject.GET("/:id", h.FindDetailSubject)
	}
}
