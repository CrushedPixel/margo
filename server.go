package margo

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"errors"
	"net/http"
)

type Server struct {
	*gin.Engine
	ErrorHandler ErrorHandlerFunc
}

type ErrorHandlerFunc func(context *gin.Context, r interface{})

var defaultErrorHandler ErrorHandlerFunc = func(c *gin.Context, r interface{}) {
	if err, ok := r.(error); ok {
		logInfo(fmt.Sprintf("Error handling request: %s\n", err.Error()))
	} else {
		logInfo(fmt.Sprintf("Error handling request: %+v\n", r))
	}

	c.Status(http.StatusInternalServerError)
}

func NewServer() *Server {
	g := gin.New()

	return &Server{
		Engine:       g,
		ErrorHandler: defaultErrorHandler,
	}
}

func (s *Server) Register(e *Endpoint) (gin.IRoutes, error) {
	logInfo(fmt.Sprintf("Registering endpoint %s", e.String()))

	if len(e.Handlers) < 1 {
		return nil, errors.New("at least one endpoint handler required")
	}

	handlers := make([]HandlerFunc, len(e.Handlers)+2)
	handlers[0] = queryParamsValidator(e.QueryParams)
	handlers[1] = bodyParamsValidator(e.BodyParams)

	for i := range e.Handlers {
		handlers[2+i] = e.Handlers[i]
	}

	return s.Handle(
		e.Method, e.Path,
		s.toGinHandler(handlers),
	), nil
}

// converts margo handlers into a single gin handler
func (s *Server) toGinHandler(handlers []HandlerFunc) gin.HandlerFunc {
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
