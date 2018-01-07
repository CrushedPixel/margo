package margo

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"errors"
	"net/http"
)

type Application struct {
	*gin.Engine
	ErrorHandler ErrorHandlerFunc
}

type ErrorHandlerFunc func(context *gin.Context, r interface{})

func defaultErrorHandler(c *gin.Context, r interface{}) {
	if err, ok := r.(error); ok {
		logInfo(fmt.Sprintf("Error handling request: %s\n", err.Error()))
	} else {
		logInfo(fmt.Sprintf("Error handling request: %+v\n", r))
	}

	c.Status(http.StatusInternalServerError)
}

func NewServer() *Application {
	g := gin.New()

	return &Application{
		Engine:       g,
		ErrorHandler: defaultErrorHandler,
	}
}

func (s *Application) Register(e Endpoint) gin.IRoutes {
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
		panic("Endpoint must not return nil")
	}
}
