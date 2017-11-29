package margo

import "github.com/gin-gonic/gin"

const (
	bodyParams  = "margoBodyParams"
	queryParams = "margoQueryParams"
)

type Context struct {
	*gin.Context
}

func (c *Context) GetBodyParams() interface{} {
	p, _ := c.Get(bodyParams)
	return p
}

func (c *Context) GetQueryParams() interface{} {
	p, _ := c.Get(queryParams)
	return p
}
