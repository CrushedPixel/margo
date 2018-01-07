// Package margo is a web framework providing a thin abstraction over the gin web framework.
package margo

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"errors"
	"net/http"
)

// ErrorHandlerFunc is a function handling any errors occurring during
// execution of an Endpoint's HandlerChain.
type ErrorHandlerFunc func(context *gin.Context, r interface{})

func defaultErrorHandler(c *gin.Context, r interface{}) {
	if err, ok := r.(error); ok {
		logInfo("Error handling request: %s\n", err.Error())
	} else {
		logInfo("Error handling request: %+v\n", r)
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
	logInfo(fmt.Sprintf("Registering endpoint %s %s", e.Method(), e.Path()))

	handlers := e.Handlers()
	if len(handlers) < 1 {
		panic(errors.New("at least one endpoint handler required"))
	}

	return s.Handle(e.Method(), e.Path(), s.toGinHandler(e.Handlers()))
}

// toGinHandler converts a HandlerChain into a single gin.HandlerFunc
func (s *Application) toGinHandler(handlers HandlerChain) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				s.ErrorHandler(c, r)
			}
		}()

		context := &Context{c}

		for _, h := range handlers {
			if response := h(context); response != nil {
				err := response.Send(c)
				if err != nil {
					panic(err)
				}
				return
			}
		}

		// if we're here, the final handler hasn't returned a value
		panic(errors.New("endpoint must not return nil"))
	}
}
