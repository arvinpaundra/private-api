package submission

import (
	"github.com/arvinpaundra/private-api/application/rest/handler"
	"github.com/arvinpaundra/private-api/core/validator"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SubmissionRouter struct {
	db  *gorm.DB
	vld *validator.Validator
}

func NewSubmissionRouter(db *gorm.DB, vld *validator.Validator) *SubmissionRouter {
	return &SubmissionRouter{
		db:  db,
		vld: vld,
	}
}

func (r *SubmissionRouter) Public(g *gin.RouterGroup) {
	h := handler.NewSubmissionHandler(r.db, r.vld)

	submission := g.Group("/modules/:module_slug/submissions")
	{
		submission.POST("", h.StartSubmission)
		submission.POST("/:submission_code/answers", h.SubmitAnswer)
		submission.PATCH("/:submission_code/finalize", h.FinalizeSubmission)
	}
}
