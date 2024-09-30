package types

import (
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Map map[string]any

type ApiError struct {
	Status int
	Msg    string
}

func (a ApiError) Error() string {
	return a.Msg
}

func NewApiError(status int, msg string) ApiError {
	return ApiError{
		Status: status,
		Msg:    msg,
	}
}

func NewPasswordError() ApiError {
	return ApiError{
		Status: http.StatusUnauthorized,
		Msg:    "Password should be minimun 8 character long, with 1 Uppercase, 1 Lowercase, 1 number and 1 special character at least",
	}
}

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}
