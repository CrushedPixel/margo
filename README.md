# margo
A small but useful wrapper library around [gin](https://github.com/gin-gonic/gin).

## Motivation
**gin** is an amazing framework, but when writing handler functions and middleware one quickly gets lost outputting 
data to the user.

**margo** solves this by having its handler functions return an object implementing the `Response` interface, 
whose `send` method is responsible for transmitting data to the client.

Additionally, **margo** comes with built-in middleware binding query and body parameters into the context, 
reducing manual parameter validation to a minimum.  

Here's a comparison:
```go
type ExampleQueryParams struct {
	// binding:"required" ensures the parameter is set (go-playground/validator.v8)
	Message string `form:"message" binding:"required"`
}

func usingMargo() {
	s := margo.NewServer()

	endpoint := margo.GetEndpoint("/", func(context *margo.Context) margo.Response {
		qp := context.GetQueryParams().(*ExampleQueryParams)

		if qp.Message == "Hello World" {
			return margo.BadRequest(margo.InvalidParamsError("message", "too uncreative"))
		}

		return margo.OK(gin.H{
			"message": qp.Message,
		})
	})

	// register ExampleQueryParams struct for binding to this endpoint
	endpoint.QueryParams = ExampleQueryParams{}

	// register endpoint to server
	s.Register(endpoint)

	s.Run("127.0.0.1:8080")
}

func usingGin() {
	s := gin.Default()

	s.GET("/", func(c *gin.Context) {
		qp := &ExampleQueryParams{}
		if err := c.BindQuery(qp); err != nil {
			// TODO: parse error object, construct some pretty error response
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "missing parameter: 'message'",
			})
			return
		}

		if qp.Message == "Hello World" {
			// TODO: create consistent error output
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "parameter 'message' is too uncreative",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": qp.Message,
		})
	})

	s.Run("127.0.0.1:8080")
}

```

Note that a `margo.Server` is merely a wrapper around `gin.Engine`, so you may use the underlying technology to its full extent.

You can also create `margo.Endpoint` structs programatically and register them to the server using `s.Register(endpoint)`,
which makes **margo** suitable for libraries that need to dynamically create endpoints.

## Outlook
I plan on using `margo` as the base for a larger web framework, therefore the API may be subject to improvements.