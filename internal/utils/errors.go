package utils

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

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

func UniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
