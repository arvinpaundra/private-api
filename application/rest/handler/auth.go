package handler

import (
	"net/http"

	"github.com/arvinpaundra/private-api/core/format"
	"github.com/arvinpaundra/private-api/core/token"
	"github.com/arvinpaundra/private-api/core/validator"
	"github.com/arvinpaundra/private-api/domain/auth/constant"
	"github.com/arvinpaundra/private-api/domain/auth/service"
	"github.com/arvinpaundra/private-api/infrastructure/auth"
	"github.com/arvinpaundra/private-api/infrastructure/shared"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db        *gorm.DB
	rdb       *redis.Client
	validator *validator.Validator
}

func NewAuthHandler(db *gorm.DB, rdb *redis.Client, validator *validator.Validator) *AuthHandler {
	return &AuthHandler{
		db:        db,
		rdb:       rdb,
		validator: validator,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var command service.UserRegisterCommand

	err := c.ShouldBindJSON(&command)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, format.UnprocessableEntity(err.Error()))
		return
	}

	verrs := h.validator.Validate(command)
	if verrs != nil {
		c.JSON(http.StatusBadRequest, format.BadRequest("invalid request body", verrs))
		return
	}

	// Create Asynq publisher for domain events
	svc := service.NewUserRegister(
		auth.NewUserReaderRepository(h.db),
		auth.NewUserWriterRepository(h.db),
		auth.NewUnitOfWork(h.db),
	)

	err = svc.Execute(c.Request.Context(), command)
	if err != nil {
		switch err {
		case constant.ErrEmailAlreadyExists, constant.ErrUsernameAlreadyExists:
			c.JSON(http.StatusConflict, format.Conflict(err.Error()))
		default:
			c.JSON(http.StatusInternalServerError, format.InternalServerError())
		}
		return
	}

	c.JSON(http.StatusCreated, format.SuccessCreated("user registered successfully", nil))
}

func (h *AuthHandler) Login(c *gin.Context) {
	var command service.UserLoginCommand

	err := c.ShouldBindJSON(&command)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, format.UnprocessableEntity(err.Error()))
		return
	}

	verrs := h.validator.Validate(command)
	if verrs != nil {
		c.JSON(http.StatusBadRequest, format.BadRequest("invalid request body", verrs))
		return
	}

	svc := service.NewUserLogin(
		auth.NewUserReaderRepository(h.db),
		auth.NewUserWriterRepository(h.db),
		token.NewJWT(viper.GetString("JWT_SECRET")),
		auth.NewUnitOfWork(h.db),
	)

	result, err := svc.Execute(c.Request.Context(), command)
	if err != nil {
		switch err {
		case constant.ErrUserNotFound:
			c.JSON(http.StatusNotFound, format.NotFound(err.Error()))
		default:
			c.JSON(http.StatusInternalServerError, format.InternalServerError())
		}
		return
	}

	c.JSON(http.StatusOK, format.SuccessOK("user logged in successfully", result))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var command service.UserLogoutCommand

	err := c.ShouldBindJSON(&command)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, format.UnprocessableEntity(err.Error()))
		return
	}

	verrs := h.validator.Validate(command)
	if verrs != nil {
		c.JSON(http.StatusBadRequest, format.BadRequest("invalid request body", verrs))
		return
	}

	svc := service.NewUserLogout(
		auth.NewUserReaderRepository(h.db),
		auth.NewUserWriterRepository(h.db),
		token.NewJWT(viper.GetString("JWT_SECRET")),
		shared.NewAuthStorage(c),
		auth.NewUnitOfWork(h.db),
	)

	err = svc.Execute(c.Request.Context(), command)
	if err != nil {
		switch err {
		case constant.ErrUserNotFound:
			c.JSON(http.StatusNotFound, format.NotFound(err.Error()))
		default:
			c.JSON(http.StatusInternalServerError, format.InternalServerError())
		}
		return
	}

	c.JSON(http.StatusOK, format.SuccessOK("user logged out successfully", nil))
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var command service.RefreshTokenCommand

	err := c.ShouldBindJSON(&command)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, format.UnprocessableEntity(err.Error()))
		return
	}

	verrs := h.validator.Validate(command)
	if verrs != nil {
		c.JSON(http.StatusBadRequest, format.BadRequest("invalid request body", verrs))
		return
	}

	svc := service.NewRefreshToken(
		auth.NewUserReaderRepository(h.db),
		auth.NewUserWriterRepository(h.db),
		token.NewJWT(viper.GetString("JWT_SECRET")),
		auth.NewUnitOfWork(h.db),
	)

	result, err := svc.Execute(c.Request.Context(), command)
	if err != nil {
		switch err {
		case constant.ErrInvalidRefreshToken:
			c.JSON(http.StatusUnauthorized, format.Unauthorized(err.Error()))
		default:
			c.JSON(http.StatusInternalServerError, format.InternalServerError())
		}
		return
	}

	c.JSON(http.StatusOK, format.SuccessOK("token refreshed successfully", result))
}
