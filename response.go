package gohttp

import "net/http"

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

// NewResponse builds a `gohttp.Response` object from an `http.Response` object.
func NewResponse(resp *http.Response) (*Response, error) {
	json, err := ParseJSON(resp.Body)
	if err != nil {
		return nil, err
	}
	return &Response{
		Code: resp.StatusCode,
		Body: json,
	}, err
}
