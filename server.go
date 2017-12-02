package margo

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"errors"
	"net/http"
)

type Server struct {
	*gin.Engine
}

func NewServer() *Server {
	g := gin.New()
	g.Use(gin.Logger())
	g.Use(gin.Recovery())

	return &Server{g}
}

func (s *Server) Register(e *Endpoint) (gin.IRoutes, error) {
	if !e.Enabled {
		logInfo(fmt.Sprintf("Disabled endpoint %s", e.String()))
		return nil, nil
	}

	logInfo(fmt.Sprintf("Registering endpoint %s", e.String()))

	if len(e.Handlers) < 1 {
		return nil, errors.New("at least one endpoint handler required")
	}

	handlers := make([]HandlerFunc, len(e.Handlers)+3)
	handlers[0] = queryParamsValidator(e.QueryParams)
	handlers[1] = bodyParamsValidator(e.BodyParams)

	for i := range e.Handlers {
		handlers[2+i] = e.Handlers[i]
	}

	return s.Handle(
		e.Method, e.Path,
		toGinHandler(handlers),
	), nil
}

// converts margo handlers into a single gin handler
func toGinHandler(handlers []HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					logInfo(fmt.Sprintf("Error handling request: %s\n", err.Error()))
				} else {
					logInfo(fmt.Sprintf("Error handling request: %+v\n", r))
				}

				res := NewErrorResponse(http.StatusInternalServerError, InternalServerError())
				res.Send(c)
			}
		}()

		context := &Context{c}

		for _, h := range handlers {
			if response := h(context); response != nil {
				response.Send(c)
				return
			}
		}

		// if we're here, the final handler hasn't returned a value
		panic("Endpoint must not return nil")
	}
}
