package margo

import (
	"github.com/satori/go.uuid"
	"net/http"
	"fmt"
	"strconv"
	"github.com/gin-gonic/gin"
)

type MargoError struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

func NewError(status int, code string, detail string) *MargoError {
	return &MargoError{
		ID:     uuid.NewV4().String(),
		Status: strconv.Itoa(status),
		Code:   code,
		Detail: detail,
	}
}

func (e *MargoError) Send(c *gin.Context) error {
	s, err := strconv.Atoi(e.Status)
	if err != nil {
		panic(err)
	}
	return NewErrorResponse(s, e).Send(c)
}

const (
	InvalidParams = "INVALID_PARAMS"
	MissingParams = "MISSING_PARAMS"
	Internal      = "INTERNAL"
)

func InvalidParamsError(field string, validation string) *MargoError {
	details := ""
	if field != "" {
		details = fmt.Sprintf("Invalid parameter: %s", field)
		if validation != "" {
			details += fmt.Sprintf(" (%s)", validation)
		}
	}

	return NewError(http.StatusBadRequest, InvalidParams, details)
}

func MissingParamsError(typ string) *MargoError {
	return NewError(http.StatusBadRequest, MissingParams, typ)
}

func InternalServerError() *MargoError {
	return NewError(http.StatusInternalServerError, Internal, "")
}
