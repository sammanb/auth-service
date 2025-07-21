package utils

type AppError struct {
	Code    int
	Message string
}

func (a *AppError) Error() string {
	return a.Message
}

func NewAppError(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}
