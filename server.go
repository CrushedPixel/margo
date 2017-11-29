package margo

import (
	"github.com/gin-gonic/gin"
)

type Context struct {
	*gin.Context
}

type Server struct {
	*gin.Engine
}

func NewServer() *Server {
	g := gin.New()
	g.Use(gin.Logger())
	g.Use(gin.Recovery())

	return &Server{g}
}

func (s *Server) Endpoint(e *Endpoint) {
	//if !e.Enabled {
	//	return // TODO: log on startup
	//}

	println("HERE0")

	s.Handle(
		e.Verb, e.Path,

		chain(
			pathParamsValidator(e.PathParams),
			queryParamsValidator(e.QueryParams),
			bodyParamsValidator(e.BodyParams),
			e.Handler,
		),
	)
}

func chain(middleware ...HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		context := &Context{c}

		for _, m := range middleware {
			if response := m(context); response != nil {
				response.Send(c)
				return
			}
		}

		// if we're here, the final handler hasn't returned a value
		panic("Endpoint must not return nil")
	}
}