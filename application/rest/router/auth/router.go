package auth

import (
	"github.com/arvinpaundra/private-api/application/rest/handler"
	"github.com/arvinpaundra/private-api/application/rest/middleware"
	"github.com/arvinpaundra/private-api/core/validator"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthRouter struct {
	db     *gorm.DB
	logger *zap.Logger
	vld    *validator.Validator
}

func NewAuthRouter(
	db *gorm.DB,
	logger *zap.Logger,
	vld *validator.Validator,
) *AuthRouter {
	return &AuthRouter{
		db:     db,
		logger: logger,
		vld:    vld,
	}
}

func (r *AuthRouter) Public(g *gin.RouterGroup) {
	h := handler.NewAuthHandler(r.db, r.logger, r.vld)

	auth := g.Group("/auth")
	{
		auth.POST("/login", h.Login)
		auth.POST("/register", h.Register)
		// auth.POST("/refresh-token", h.RefreshToken)
	}
}

func (r *AuthRouter) Private(g *gin.RouterGroup) {
	h := handler.NewAuthHandler(r.db, r.logger, r.vld)
	m := middleware.NewAuthenticate(r.db)

	me := g.Group("/me", m.Authenticate())
	{
		me.POST("/logout", h.Logout)
	}
}
