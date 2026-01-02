package handler

import (
	"net/http"

	"github.com/arvinpaundra/private-api/core/format"
	"github.com/arvinpaundra/private-api/domain/dashboard/service"
	"github.com/arvinpaundra/private-api/infrastructure/dashboard"
	"github.com/arvinpaundra/private-api/infrastructure/shared"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DashboardHandler struct {
	db *gorm.DB
}

func NewDashboardHandler(db *gorm.DB) *DashboardHandler {
	return &DashboardHandler{
		db: db,
	}
}

func (h *DashboardHandler) GetStatistics(c *gin.Context) {
	svc := service.NewGetDashboardStatistics(
		shared.NewAuthStorage(c),
		dashboard.NewModuleACLAdapter(h.db),
		dashboard.NewSubjectACLAdapter(h.db),
		dashboard.NewGradeACLAdapter(h.db),
		dashboard.NewSubmissionACLAdapter(h.db),
		dashboard.NewUserACLAdapter(h.db),
	)

	result, err := svc.Execute(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, format.InternalServerError())
		return
	}

	c.JSON(http.StatusOK, format.SuccessOK("statistics retrieved successfully", result))
}
