/*
Package margo is a web framework providing a thin abstraction over the gin web framework.

It introduces the concept of handler functions returning Response values instead of directly
setting headers or writing data to the response body.
This makes it clear where handler functions write data, thus greatly improving code quality.

Basic Example:

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
*/
package margo
