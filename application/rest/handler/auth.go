package handler

import (
	"net/http"

	"github.com/arvinpaundra/private-api/config"
	"github.com/arvinpaundra/private-api/core/format"
	"github.com/arvinpaundra/private-api/core/token"
	"github.com/arvinpaundra/private-api/core/validator"
	"github.com/arvinpaundra/private-api/domain/auth/constant"
	"github.com/arvinpaundra/private-api/domain/auth/service"
	"github.com/arvinpaundra/private-api/infrastructure/auth"
	"github.com/arvinpaundra/private-api/infrastructure/shared"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db        *gorm.DB
	logger    *zap.Logger
	validator *validator.Validator
}

func NewAuthHandler(
	db *gorm.DB,
	logger *zap.Logger,
	validator *validator.Validator,
) *AuthHandler {
	return &AuthHandler{
		db:        db,
		logger:    logger.With(zap.String("domain", "auth")),
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
		h.logger.Error("failed to register user", zap.Error(err))

		switch err {
		case constant.ErrEmailAlreadyExists:
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
		token.NewJWT(config.GetString("JWT_SECRET")),
		auth.NewUnitOfWork(h.db),
	)

	result, err := svc.Execute(c.Request.Context(), command)
	if err != nil {
		h.logger.Error("failed to login user", zap.Error(err))

		switch err {
		case constant.ErrUserNotFound:
			c.JSON(http.StatusNotFound, format.NotFound(err.Error()))
		case constant.ErrWrongEmailOrPassword:
			c.JSON(http.StatusUnauthorized, format.Unauthorized(err.Error()))
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
		token.NewJWT(config.GetString("JWT_SECRET")),
		shared.NewAuthStorage(c),
		auth.NewUnitOfWork(h.db),
	)

	err = svc.Execute(c.Request.Context(), command)
	if err != nil {
		h.logger.Error("failed to logout user", zap.Error(err))

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
		token.NewJWT(config.GetString("JWT_SECRET")),
		auth.NewUnitOfWork(h.db),
	)

	result, err := svc.Execute(c.Request.Context(), command)
	if err != nil {
		h.logger.Error("failed to refresh token", zap.Error(err))

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
