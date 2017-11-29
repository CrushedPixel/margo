package margo

import (
	"github.com/satori/go.uuid"
	"net/http"
	"fmt"
)

type MargoError struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

func newError(status int, code string, detail string) *MargoError {
	return &MargoError{
		ID:     uuid.NewV4().String(),
		Status: string(status),
		Code:   code,
		Detail: detail,
	}
}

const (
	InvalidParams = "INVALID_PARAMS"
)

func InvalidParamsError(field *string) *MargoError {
	var details string
	if field != nil {
		details = fmt.Sprintf("Invalid parameter: %s", field)
	} else {
		details = "Missing parameters"
	}

	return newError(http.StatusBadRequest, InvalidParams, details)
}
