package margo

type Endpoint struct {
	Enabled bool

	Verb string
	Path string

	PathParams  interface{}
	QueryParams interface{}
	BodyParams  interface{}

	Handler HandlerFunc
}

type HandlerFunc func(context *Context) Response