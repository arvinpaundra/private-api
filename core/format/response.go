package format

import (
	"net/http"

	"github.com/arvinpaundra/private-api/core/validator"
)

type Meta struct {
	Code    int    `json:"code"`
	Message string `json:"message"`

	Pagination *Pagination `json:"pagination,omitempty"`
}

type Response struct {
	Meta   Meta            `json:"meta"`
	Data   any             `json:"data"`
	Errors validator.Error `json:"errors,omitempty"`
}

// 200 - OK
func SuccessOK(message string, data any, pagination ...Pagination) Response {
	var p *Pagination

	if len(pagination) > 0 {
		p = &pagination[0]
	}

	return Response{
		Meta: Meta{
			Code:       http.StatusOK,
			Message:    message,
			Pagination: p,
		},
		Data: data,
	}
}

// 201 - Created
func SuccessCreated(message string, data any) Response {
	return Response{
		Meta: Meta{
			Code:    http.StatusCreated,
			Message: message,
		},
		Data: data,
	}
}

// 400 - Bad Request
func BadRequest(message string, errors validator.Error) Response {
	return Response{
		Meta: Meta{
			Code:    http.StatusBadRequest,
			Message: message,
		},
		Errors: errors,
	}
}

// 401 - Unauthorized
func Unauthorized(message string) Response {
	return Response{
		Meta: Meta{
			Code:    http.StatusUnauthorized,
			Message: message,
		},
	}
}

// 403 - Forbidden
func Forbidden(message string) Response {
	return Response{
		Meta: Meta{
			Code:    http.StatusForbidden,
			Message: message,
		},
	}
}

// 404 - Not Found
func NotFound(message string) Response {
	return Response{
		Meta: Meta{
			Code:    http.StatusNotFound,
			Message: message,
		},
	}
}

// 409 - Conflict
func Conflict(message string) Response {
	return Response{
		Meta: Meta{
			Code:    http.StatusConflict,
			Message: message,
		},
	}
}

// 422 - Unprocessable Entity
func UnprocessableEntity(message string) Response {
	return Response{
		Meta: Meta{
			Code:    http.StatusUnprocessableEntity,
			Message: message,
		},
	}
}

// 500 - Internal Server Error
func InternalServerError() Response {
	return Response{
		Meta: Meta{
			Code:    http.StatusInternalServerError,
			Message: "internal server error",
		},
	}
}
