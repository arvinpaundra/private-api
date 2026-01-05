package router

import (
	"github.com/arvinpaundra/private-api/application/rest/middleware"
	"github.com/arvinpaundra/private-api/application/rest/router/auth"
	"github.com/arvinpaundra/private-api/application/rest/router/dashboard"
	"github.com/arvinpaundra/private-api/application/rest/router/grade"
	"github.com/arvinpaundra/private-api/application/rest/router/health"
	"github.com/arvinpaundra/private-api/application/rest/router/module"
	"github.com/arvinpaundra/private-api/application/rest/router/subject"
	"github.com/arvinpaundra/private-api/application/rest/router/submission"
	"github.com/arvinpaundra/private-api/core/validator"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(g *gin.Engine, db *gorm.DB) *gin.Engine {
	g.Use(middleware.Cors())
	g.Use(gin.Recovery())
	g.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/metrics"},
	}))

	// Health check endpoints
	healthRouter := health.NewHealthRouter(db)
	healthRouter.Register(g)

	v1 := g.Group("/v1")

	authRouter := auth.NewAuthRouter(db, validator.NewValidator())
	subjectRouter := subject.NewSubjectRouter(db, validator.NewValidator())
	gradeRouter := grade.NewGradeRouter(db, validator.NewValidator())
	moduleRouter := module.NewModuleRouter(db, validator.NewValidator())
	submissionRouter := submission.NewSubmissionRouter(db, validator.NewValidator())
	dashboardRouter := dashboard.NewDashboardRouter(db)

	// public routes
	authRouter.Public(v1)
	moduleRouter.Public(v1)
	submissionRouter.Public(v1)

	// private routes
	authRouter.Private(v1)
	subjectRouter.Private(v1)
	gradeRouter.Private(v1)
	moduleRouter.Private(v1)
	submissionRouter.Private(v1)
	dashboardRouter.Private(v1)

	return g
}
