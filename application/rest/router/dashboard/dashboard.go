package dashboard

import (
	"github.com/arvinpaundra/private-api/application/rest/handler"
	"github.com/arvinpaundra/private-api/application/rest/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DashboardRouter struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewDashboardRouter(
	db *gorm.DB,
	logger *zap.Logger,
) *DashboardRouter {
	return &DashboardRouter{
		db:     db,
		logger: logger,
	}
}

func (r *DashboardRouter) Private(g *gin.RouterGroup) {
	h := handler.NewDashboardHandler(r.db, r.logger)
	m := middleware.NewAuthenticate(r.db)

	g.GET("/dashboard/statistics", m.Authenticate(), h.GetStatistics)
}
