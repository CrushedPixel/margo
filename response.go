package margo

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response interface {
	Send(context *gin.Context)
}

type ErrorResponse struct {
	Status int
	Errors []*MargoError
}

func (r *ErrorResponse) Send(c *gin.Context) {
	c.JSON(r.Status, gin.H{"errors": r.Errors})
}

type DataResponse struct {
	Status int
	Data   interface{}
}

func (r *DataResponse) Send(c *gin.Context) {
	c.JSON(r.Status, gin.H{"data": r.Data})
}

func BadRequest(error []*MargoError) *ErrorResponse {
	return &ErrorResponse{
		Status: http.StatusBadRequest,
		Errors: error,
	}
}

func OK(data interface{}) *DataResponse {
	return &DataResponse{
		Status: http.StatusOK,
		Data: data,
	}
}