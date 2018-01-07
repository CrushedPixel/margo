package margo

import "github.com/gin-gonic/gin"

const (
	bodyParams  = "__margoBodyParams"
	queryParams = "__margoQueryParams"
)

// A Context embeds a *gin.Context, adding
// getter methods for parsed body and query parameters.
type Context struct {
	*gin.Context
}

// BodyParams returns a pointer to the model instance bound by a BindingEndpoint.
// Returns nil if no body parameter binding was done.
func (c *Context) BodyParams() interface{} {
	p, _ := c.Get(bodyParams)
	return p
}

// QueryParams returns a pointer to the model instance bound by a BindingEndpoint.
// Returns nil if no query parameter binding was done.
func (c *Context) QueryParams() interface{} {
	p, _ := c.Get(queryParams)
	return p
}
