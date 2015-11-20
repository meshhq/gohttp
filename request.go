package gohttp

// Request models an HTTP request that can be sumbitted to a `Client` for
// executing.
type Request struct {

	// Method models the HTTP method to be used for the request.
	Method string

	// URL models the API route to be used for the request.
	URL string

	// Bodel models the request body to be used for the request.
	Body map[string]interface{}
}
