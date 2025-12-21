package grade

import (
	"github.com/arvinpaundra/private-api/application/rest/handler"
	"github.com/arvinpaundra/private-api/application/rest/middleware"
	"github.com/arvinpaundra/private-api/core/validator"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GradeRouter struct {
	db  *gorm.DB
	vld *validator.Validator
}

func NewGradeRouter(db *gorm.DB, vld *validator.Validator) *GradeRouter {
	return &GradeRouter{
		db:  db,
		vld: vld,
	}
}

func (r *GradeRouter) Private(g *gin.RouterGroup) {
	h := handler.NewGradeHandler(r.db, r.vld)
	m := middleware.NewAuthenticate(r.db)

	grade := g.Group("/grades", m.Authenticate())
	{
		grade.POST("", h.CreateGrade)
		grade.PUT("/:id", h.UpdateGrade)
		grade.DELETE("/:id", h.DeleteGrade)
		grade.GET("", h.FindAllGrades)
		grade.GET("/:id", h.FindDetailGrade)
	}
}
