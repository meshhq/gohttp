package gohttp

// Request models an HTTP request that can be sumbitted to a `Client` for
// executing.
type Request struct {

	// Method models the HTTP method to be used for the request.
	Method string

	// URL models the API route to be used for the request.
	URL string

	// Params contains the URL parameters to be used with the request.
	Params map[string]string

	// Body contains the request body to be used for the request.
	Body interface{}
}
