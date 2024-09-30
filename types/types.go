package types

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
