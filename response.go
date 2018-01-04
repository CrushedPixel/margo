package margo

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"io"
	"strconv"
)

const contentLengthHeader  = "Content-Length"

// base Response interface
type Response interface {
	Send(context *gin.Context) error
}

type doNothingResponse struct{}

func (r *doNothingResponse) Send(c *gin.Context) error {
	return nil
}

// Generic JSON Response
type jsonResponse struct {
	Status int
	Data   interface{}
}

func (r *jsonResponse) Send(c *gin.Context) error {
	c.JSON(r.Status, r.Data)
	return nil
}

// Empty Response only setting a status code
type emptyResponse struct {
	Status int
}

func (r *emptyResponse) Send(c *gin.Context) error {
	c.Status(r.Status)
	c.Writer.WriteHeaderNow()
	return nil
}

// Response sending a file
type fileResponse struct {
	file *os.File
}

func (r *fileResponse) Send(c *gin.Context) error {
	stat, err := r.file.Stat()
	if err != nil {
		return err
	}

	c.Status(http.StatusOK)
	c.Header(contentLengthHeader, strconv.FormatInt(stat.Size(), 10))

	// ignore any errors writing the file,
	// as they are most likely caused by
	// the client closing the connection
	io.Copy(c.Writer, r.file)
	return nil
}

// Utility methods to create responses
func JSON(status int, data interface{}) *jsonResponse {
	return &jsonResponse{
		status, data,
	}
}

func JSON200(data interface{}) Response {
	return JSON(http.StatusOK, data)
}

func NewEmptyResponse(status int) Response {
	return &emptyResponse{status}
}

func NewErrorResponse(status int, errors ...*MargoError) Response {
	return JSON(status, gin.H{"errors": errors})
}

func DoNothing() Response {
	return &doNothingResponse{}
}

func BadRequest(error ...*MargoError) Response {
	return NewErrorResponse(http.StatusBadRequest, error...)
}

func SendFile(file *os.File) Response {
	return &fileResponse{
		file: file,
	}
}
