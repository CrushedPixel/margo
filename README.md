# margo
[![GoDoc](https://godoc.org/github.com/CrushedPixel/margo?status.svg)](https://godoc.org/github.com/CrushedPixel/margo) [![Go Report Card](https://goreportcard.com/badge/github.com/crushedpixel/margo)](https://goreportcard.com/report/github.com/crushedpixel/margo)

A tiny framework on top of [gin](https://github.com/gin-gonic/gin).

## Motivation
**gin** is an amazing framework, but when writing handler functions and middleware it is easy to get confused
keeping track where headers are set and data is written to the response.

**margo** solves this by having its handler functions return an object implementing the `Response` interface, 
whose `Send` method is responsible for transmitting data to the client.

## Example usage
```go
type errorResponse struct {
    err error
}

// satisfies margo.Response
func (r *errorResponse) Send(context *gin.Context) error {
    context.String(http.StatusInternalServerError, "an internal server error occurred: %s", r.err.Error())
    return nil
}

func newErrorResponse(err error) margo.Response {
    return &errorResponse{
        err: err,
    }
}

func main() {
    // create new Application instance
    app := margo.NewApplication()

    // create endpoint handling the index route
    endpoint := margo.GET("/", func(context *gin.Context) margo.Response {
        // handle the request however you wish, for example
        // parse some request parameters
        params, err := parseQueryParams(context)
        if err != nil {
            // handle the error however you like,
            // for example by returning a generic error response
            // which implements margo.Response
            return newErrorResponse(err)
        }

        // do something with the retrieved params,
        // for example output them as json using
        // the builtin JSON function
        return margo.JSON(http.StatusOK, params)
    })
    app.Endpoint(endpoint) // register endpoint with Application

    app.Run(":8080") // run application
}
```

Note that a `margo.Application` is merely a wrapper around `gin.Engine`, so you may use the underlying technology to its full extent.

## Additional libraries
- [binder](https://github.com/CrushedPixel/binder) is a request binding middleware for margo.
- [jargo](https://github.com/CrushedPixel/jargo) is a fully-featured jsonapi web framework built on top of margo.
