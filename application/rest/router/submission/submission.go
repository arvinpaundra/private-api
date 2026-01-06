package submission

import (
	"github.com/arvinpaundra/private-api/application/rest/handler"
	"github.com/arvinpaundra/private-api/application/rest/middleware"
	"github.com/arvinpaundra/private-api/core/validator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SubmissionRouter struct {
	db     *gorm.DB
	logger *zap.Logger
	vld    *validator.Validator
}

func NewSubmissionRouter(
	db *gorm.DB,
	logger *zap.Logger,
	vld *validator.Validator,
) *SubmissionRouter {
	return &SubmissionRouter{
		db:     db,
		logger: logger,
		vld:    vld,
	}
}

func (r *SubmissionRouter) Private(g *gin.RouterGroup) {
	h := handler.NewSubmissionHandler(r.db, r.logger, r.vld)
	m := middleware.NewAuthenticate(r.db)

	submission := g.Group("/submissions", m.Authenticate())

	submission.GET("", h.GetAllSubmissions)
}

func (r *SubmissionRouter) Public(g *gin.RouterGroup) {
	h := handler.NewSubmissionHandler(r.db, r.logger, r.vld)

	submission := g.Group("/modules/:module_slug/submissions")
	{
		submission.POST("", h.StartSubmission)
		submission.POST("/:submission_code/answers", h.SubmitAnswer)
		submission.PATCH("/:submission_code/finalize", h.FinalizeSubmission)
	}
}
