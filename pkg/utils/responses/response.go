package responses

import (
	"github.com/shlembo598/text-lexicon-go/pkg/utils/httpErrors"
)

const (
	statusSuccess = "success"
	statusError   = "error"
)

const (
	ErrBadRequest         = "Bad request"
	ErrEmailAlreadyExists = "User with given email already exists"
	ErrNoSuchUser         = "User not found"
	ErrWrongCredentials   = "Wrong Credentials"
	ErrNotFound           = "Not Found"
	ErrUnauthorized       = "Unauthorized"
	ErrForbidden          = "Forbidden"
	ErrBadQueryParams     = "Invalid query params"
)

type Response struct {
	Status string             `json:"status"`
	Data   interface{}        `json:"data,omitempty"`
	Error  httpErrors.RestErr `json:"error,omitempty"`
}

type APIError struct {
	Status int         `json:"code"`
	Error  interface{} `json:"error"`
}

func SuccessResponse(data interface{}) interface{} {
	return &Response{
		Status: statusSuccess,
		Data:   data,
		Error:  nil,
	}
}

func ErrorResponse(err error) (int, interface{}) {
	response := &Response{
		Status: statusError,
		Data:   nil,
		Error:  httpErrors.ParseErrors(err),
	}
	return httpErrors.ParseErrors(err).Status(), response
}
