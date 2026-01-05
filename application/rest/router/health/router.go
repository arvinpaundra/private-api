package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HealthRouter struct {
	db *gorm.DB
}

func NewHealthRouter(db *gorm.DB) *HealthRouter {
	return &HealthRouter{
		db: db,
	}
}

func (r *HealthRouter) Register(g *gin.Engine) {
	g.GET("/livez", r.livez)
	g.GET("/readyz", r.readyz)
}

func (r *HealthRouter) livez(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (r *HealthRouter) readyz(c *gin.Context) {
	// Check database connection
	sqlDB, err := r.db.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unavailable",
			"error":  "database connection failed",
		})
		return
	}

	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "unavailable",
			"error":  "database ping failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
	})
}
