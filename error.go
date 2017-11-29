package margo

import (
	"github.com/satori/go.uuid"
	"net/http"
	"fmt"
	"strconv"
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
		Status: strconv.Itoa(status),
		Code:   code,
		Detail: detail,
	}
}

const (
	InvalidParams = "INVALID_PARAMS"
	Internal      = "INTERNAL"
)

func InvalidParamsError(field *string, validation *string) *MargoError {
	var details string
	if field != nil && validation != nil {
		details = fmt.Sprintf("Invalid parameter: %s (%s)", *field, *validation)
	} else {
		details = "Missing parameters"
	}

	return newError(http.StatusBadRequest, InvalidParams, details)
}

func InternalServerError() *MargoError {
	return newError(http.StatusInternalServerError, Internal, "")
}
