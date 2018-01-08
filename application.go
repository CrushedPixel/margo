package margo

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// ErrorHandlerFunc is a function handling any errors occurring during
// execution of an Endpoint's HandlerChain.
type ErrorHandlerFunc func(context *gin.Context, r interface{})

func defaultErrorHandler(c *gin.Context, r interface{}) {
	if err, ok := r.(error); ok {
		println(fmt.Sprintf("Error handling request: %s\n", err.Error()))
	} else {
		println(fmt.Sprintf("Error handling request: %+v\n", r))
	}

	c.Status(http.StatusInternalServerError)
}

// An Application is a thin wrapper around a gin.Engine,
// providing additional utility methods.
type Application struct {
	*gin.Engine
	// ErrorHandler is the ErrorHandlerFunc called when
	// a HandlerFunc in an Endpoint's HandlerChain
	// panics, or sending a Response returns an error.
	ErrorHandler ErrorHandlerFunc
}

// NewApplication returns a new Application with
// the underlying gin.Engine being initialized using gin.New()
// and the default error handler.
func NewApplication() *Application {
	return &Application{
		Engine:       gin.New(),
		ErrorHandler: defaultErrorHandler,
	}
}

// Endpoint exposes an Endpoint via HTTP.
func (s *Application) Endpoint(e Endpoint) gin.IRoutes {
	handlers := e.Handlers()
	if len(handlers) < 1 {
		panic(errors.New("at least one endpoint handler required"))
	}

	return s.Handle(e.Method(), e.Path(), e.Handlers().ToGinHandler(s.ErrorHandler))
}
