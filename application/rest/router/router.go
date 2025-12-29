package router

import (
	"github.com/arvinpaundra/private-api/application/rest/middleware"
	"github.com/arvinpaundra/private-api/application/rest/router/auth"
	"github.com/arvinpaundra/private-api/application/rest/router/grade"
	"github.com/arvinpaundra/private-api/application/rest/router/module"
	"github.com/arvinpaundra/private-api/application/rest/router/subject"
	"github.com/arvinpaundra/private-api/core/validator"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Register(g *gin.Engine, rdb *redis.Client, db *gorm.DB) {
	g.Use(middleware.Cors())
	g.Use(gin.Recovery())
	g.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/metrics"},
	}))

	v1 := g.Group("/v1")

	authRouter := auth.NewAuthRouter(db, rdb, validator.NewValidator())
	subjectRouter := subject.NewSubjectRouter(db, validator.NewValidator())
	gradeRouter := grade.NewGradeRouter(db, validator.NewValidator())
	moduleRouter := module.NewModuleRouter(db, validator.NewValidator())

	// public routes
	authRouter.Public(v1)

	// private routes
	authRouter.Private(v1)
	subjectRouter.Private(v1)
	gradeRouter.Private(v1)
	moduleRouter.Private(v1)
}
