package middleware

import (
	"net/http"
	"strings"

	"github.com/arvinpaundra/private-api/core/format"
	"github.com/arvinpaundra/private-api/core/token"
	"github.com/arvinpaundra/private-api/domain/auth/service"
	"github.com/arvinpaundra/private-api/domain/shared/entity"
	"github.com/arvinpaundra/private-api/infrastructure/auth"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Authenticate struct {
	db *gorm.DB
}

func NewAuthenticate(db *gorm.DB) *Authenticate {
	return &Authenticate{db: db}
}

func (a *Authenticate) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader("Authorization")

		if bearerToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, format.Unauthorized("bearer token is missing"))
			return
		}

		accessToken := strings.Replace(bearerToken, "Bearer ", "", 1)

		svc := service.NewCheckSession(
			auth.NewUserReaderRepository(a.db),
			token.NewJWT(viper.GetString("JWT_SECRET")),
		)

		command := service.CheckSessionCommand{
			AccessToken: accessToken,
		}

		res, err := svc.Execute(c.Request.Context(), command)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, format.Unauthorized(err.Error()))
			return
		}

		session := &entity.UserAuth{
			UserId:   res.UserId,
			Email:    res.Email,
			Fullname: res.Fullname,
		}

		c.Set("session", session)

		c.Next()
	}
}
