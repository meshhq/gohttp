package gohttp

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
)

// Response models a response from HTTP request.
//
// The GoHTTP Response object provides a convenient model for handling the
// result of an HTTP request.
type Response struct {

	// Code is the response code for the request.
	Code int

	// Data contains the raw response data from an request.
	Data []byte

	// Body is the deserialized response body returned from an request.
	Body interface{}

	// Error is an error that may have occured during the request.
	Error error

	// Request is the gohttp.Request object used to generate the response.
	Request *Request
}

// NewResponse builds a `gohttp.Response` object from an `http.Response` object.
func NewResponse(resp *http.Response) (*Response, error) {
	bodyContent, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var body interface{}
	if len(bodyContent) > 0 {
		canonicalHeaderKey := http.CanonicalHeaderKey(ContentType)
		for _, headerVal := range resp.Header[canonicalHeaderKey] {
			if strings.Contains(headerVal, "application/json") {
				body, err = ParseJSON(bytes.NewReader(bodyContent))
				if err != nil {
					return nil, err
				}
				break
			}
		}
		// Set body w/ no conversion if is still nil
		if body == nil {
			body = resp.Body
		}
	}

	return &Response{
		Code: resp.StatusCode,
		Body: body,
		Data: bodyContent,
	}, nil
}
