package margo

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// An Endpoint represents an HTTP endpoint
// that can be registered to an Application.
type Endpoint interface {
	// Method returns the Endpoint's HTTP method.
	Method() string
	// Path returns the Endpoint's URL.
	Path() string
	// Handlers returns a slice of handler functions
	// to be executed in order when the Endpoint is called.
	Handlers() HandlerChain
}

// basicEndpoint is a basic implementation of Endpoint.
type basicEndpoint struct {
	method   string
	path     string
	handlers HandlerChain
}

func (e *basicEndpoint) Method() string {
	return e.method
}

func (e *basicEndpoint) Path() string {
	return e.path
}

func (e *basicEndpoint) Handlers() HandlerChain {
	return e.handlers
}

// NewEndpoint returns a new Endpoint for a given HTTP method and URL path,
// with at least one HandlerFunc to be executed when the Endpoint is called.
//
// Panics if no HandlerFunc is provided.
func NewEndpoint(method string, path string, handlers ...HandlerFunc) Endpoint {
	if len(handlers) < 1 {
		panic(errors.New("at least one handler function has to be provided"))
	}

	return &basicEndpoint{
		method:   method,
		path:     path,
		handlers: HandlerChain(handlers),
	}
}

// GET returns a new GET Endpoint for a path and at least one HandlerFunc.
func GET(path string, handlers ...HandlerFunc) Endpoint {
	return NewEndpoint(http.MethodGet, path, handlers...)
}

// POST returns a new POST Endpoint for a path and at least one HandlerFunc.
func POST(path string, handlers ...HandlerFunc) Endpoint {
	return NewEndpoint(http.MethodPost, path, handlers...)
}

// PUT returns a new PUT Endpoint for a path and at least one HandlerFunc.
func PUT(path string, handlers ...HandlerFunc) Endpoint {
	return NewEndpoint(http.MethodPut, path, handlers...)
}

// PATCH returns a new PATCH Endpoint for a path and at least one HandlerFunc.
func PATCH(path string, handlers ...HandlerFunc) Endpoint {
	return NewEndpoint(http.MethodPatch, path, handlers...)
}

// DELETE returns a new DELETE Endpoint for a path and at least one HandlerFunc.
func DELETE(path string, handlers ...HandlerFunc) Endpoint {
	return NewEndpoint(http.MethodDelete, path, handlers...)
}

// A HandlerFunc is a function to be called when an Endpoint is accessed.
type HandlerFunc func(context *gin.Context) Response

// A HandlerChain is a slice of handler functions to be executed in order.
// If a HandlerFunc returns a Response value, the Response is sent to the client,
// otherwise, the next HandlerFunc in the chain is executed.
// The last HandlerFunc in the chain is expected to return a Response value.
type HandlerChain []HandlerFunc

// ToGinHandler converts a HandlerChain into a single gin.HandlerFunc.
// If any of the chain's handlers panic or Response.Send returns an error,
// the error handler is invoked.
func (chain HandlerChain) ToGinHandler(errorHandler ErrorHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				errorHandler(c, r)
			}
		}()

		for _, h := range chain {
			if response := h(c); response != nil {
				err := response.Send(c)
				if err != nil {
					panic(err)
				}
				return
			}
		}

		// if we're here, the final handler hasn't returned a value
		panic(errors.New("last handler in chain must not return nil"))
	}
}
