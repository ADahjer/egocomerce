package types

import "github.com/go-playground/validator/v10"

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

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}
