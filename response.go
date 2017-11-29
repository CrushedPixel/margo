package margo

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// base Response interface
type Response interface {
	Send(context *gin.Context)
}

// Generic JSON Response
type JSONResponse struct {
	Status int
	Data   interface{}
}

func (r *JSONResponse) Send(c *gin.Context) {
	c.JSON(r.Status, r.Data)
}

// Utility methods to create responses
func NewJSONResponse(status int, data interface{}) *JSONResponse {
	return &JSONResponse{
		status, data,
	}
}

func NewErrorResponse(status int, errors ...*MargoError) *JSONResponse {
	return NewJSONResponse(status, gin.H{"errors": errors})
}

func BadRequest(error ...*MargoError) *JSONResponse {
	return NewErrorResponse(http.StatusBadRequest, error...)
}

func OK(data interface{}) *JSONResponse {
	return NewJSONResponse(http.StatusOK, data)
}
