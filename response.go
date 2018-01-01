package margo

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// base Response interface
type Response interface {
	Send(context *gin.Context) error
}

type DoNothingResponse struct{}

func (r *DoNothingResponse) Send(c *gin.Context) error {
	return nil
}

// Generic JSON Response
type JSONResponse struct {
	Status int
	Data   interface{}
}

func (r *JSONResponse) Send(c *gin.Context) error {
	c.JSON(r.Status, r.Data)
	return nil
}

// Empty Response only setting a status code
type EmptyResponse struct {
	Status int
}

func (r *EmptyResponse) Send(c *gin.Context) error {
	c.Status(r.Status)
	c.Writer.WriteHeaderNow()
	return nil
}

// Utility methods to create responses
func NewJSONResponse(status int, data interface{}) *JSONResponse {
	return &JSONResponse{
		status, data,
	}
}

func NewEmptyResponse(status int) *EmptyResponse {
	return &EmptyResponse{status}
}

func NewErrorResponse(status int, errors ...*MargoError) *JSONResponse {
	return NewJSONResponse(status, gin.H{"errors": errors})
}

func DoNothing() *DoNothingResponse {
	return &DoNothingResponse{}
}

func BadRequest(error ...*MargoError) *JSONResponse {
	return NewErrorResponse(http.StatusBadRequest, error...)
}

func OK(data interface{}) *JSONResponse {
	return NewJSONResponse(http.StatusOK, data)
}
