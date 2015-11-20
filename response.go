package gohttp

// Response models the response from HTTP request.
type Response struct {

	// Code represents the response code for the HTTP request.
	Code int

	// Data contains the raw response data from an HTTP request.
	Data []byte

	// Body is the deserialized response body returned from an HTTP request.
	Body map[string]interface{}

	// Error represents any error that may have occured during processing.
	Error error
}
