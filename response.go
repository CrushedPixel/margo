package margo

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"io"
	"strconv"
)

// A Response is responsible for sending data to an HTTP client.
//
// When using margo, all writing of HTTP headers or content should
// happen in a Response's Send method, nowhere else in your app.
type Response interface {
	// Send sends response data to an HTTP client via a gin.Context.
	// Any errors returned are handled by the Application's ErrorHandler.
	Send(context *gin.Context) error
}

type jsonResponse struct {
	Status int
	Data   interface{}
}

func (r *jsonResponse) Send(c *gin.Context) error {
	c.JSON(r.Status, r.Data)
	return nil
}

type emptyResponse struct {
	Status int
}

func (r *emptyResponse) Send(c *gin.Context) error {
	c.Status(r.Status)
	c.Writer.WriteHeaderNow()
	return nil
}

type fileResponse struct {
	file *os.File
}

func (r *fileResponse) Send(c *gin.Context) error {
	defer r.file.Close()
	stat, err := r.file.Stat()
	if err != nil {
		return err
	}

	c.Status(http.StatusOK)
	c.Header("Content-Length", strconv.FormatInt(stat.Size(), 10))

	// ignore any errors writing the file,
	// as they are most likely caused by
	// the client closing the connection
	io.Copy(c.Writer, r.file)
	return nil
}

// JSON returns a Response sending json-encoded data
// with the specified status code.
func JSON(status int, data interface{}) Response {
	return &jsonResponse{
		status, data,
	}
}

// JSON200 returns a Response sending json-encoded data
// with status code 200 OK.
func JSON200(data interface{}) Response {
	return JSON(http.StatusOK, data)
}

// Empty returns a Response sending no data
// with the specified status code.
func Empty(status int) Response {
	return &emptyResponse{status}
}

// SendFile returns a Response sending a file
// with status code 200 OK.
func SendFile(file *os.File) Response {
	return &fileResponse{
		file: file,
	}
}
