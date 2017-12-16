package margo

import (
	"fmt"
	"net/http"
)

type HandlerFunc func(context *Context) Response

type Endpoint struct {
	Enabled bool

	Method string
	Path   string

	QueryParams interface{}
	BodyParams  interface{}

	Handlers []HandlerFunc
}

func (e *Endpoint) String() string {
	return fmt.Sprintf("%s %s", e.Method, e.Path)
}

func NewEndpoint(method string, path string, handlers ...HandlerFunc) (*Endpoint) {
	return &Endpoint{
		Method:   method,
		Path:     path,
		Handlers: handlers,
	}
}

func GetEndpoint(path string, handlers ...HandlerFunc) (*Endpoint) {
	return NewEndpoint(http.MethodGet, path, handlers...)
}

func PostEndpoint(path string, handlers ...HandlerFunc) (*Endpoint) {
	return NewEndpoint(http.MethodPost, path, handlers...)
}

func PutEndpoint(path string, handlers ...HandlerFunc) (*Endpoint) {
	return NewEndpoint(http.MethodPut, path, handlers...)
}

func PatchEndpoint(path string, handlers ...HandlerFunc) (*Endpoint) {
	return NewEndpoint(http.MethodPatch, path, handlers...)
}

func DeleteEndpoint(path string, handlers ...HandlerFunc) (*Endpoint) {
	return NewEndpoint(http.MethodDelete, path, handlers...)
}
