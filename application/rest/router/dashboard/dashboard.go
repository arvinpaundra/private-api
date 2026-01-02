package dashboard

import (
	"github.com/arvinpaundra/private-api/application/rest/handler"
	"github.com/arvinpaundra/private-api/application/rest/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DashboardRouter struct {
	db *gorm.DB
}

func NewDashboardRouter(db *gorm.DB) *DashboardRouter {
	return &DashboardRouter{
		db: db,
	}
}

func (r *DashboardRouter) Private(g *gin.RouterGroup) {
	h := handler.NewDashboardHandler(r.db)
	m := middleware.NewAuthenticate(r.db)

	g.GET("/dashboard/statistics", m.Authenticate(), h.GetStatistics)
}
