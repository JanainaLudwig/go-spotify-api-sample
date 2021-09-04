package core

import (
	"net/http"
	"runtime"
)

type AppError struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
	Status  int    `json:"status"`
	File    string `json:"-"`
	Line    int    `json:"-"`
}

func NewAppError(err AppError) *AppError {
	_, file, line, _ := runtime.Caller(1)

	err.File = file
	err.Line = line

	if err.Status == 0 {
		err.Status = http.StatusInternalServerError
	}

	return &err
}

func (a *AppError) Error() string {
	return a.Message
}

