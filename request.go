package gohttp

import (
	"net/http"
	"net/url"
)

// Request models an HTTP request.
//
// GoHTTP Request objects are sumbitted to a `gohttp.Client` for execution.
type Request struct {

	// Method is the HTTP method to be used for the request.
	Method string

	// Header are the headers for the request.
	Header http.Header

	// URL is the route to be used for the request. The URL should be relative to the BaseURL of the client.
	URL string

	// Params contains the URL parameters to be used with the request.
	Params map[string]string

	// Body contains the body to be used for the request.
	Body interface{}

	// Data contains the form data to be used for the request. Data will be sent as application/x-www-form-urlencoded.
	Data interface{}
}

//------------------------------------------------------------------------------
// Request Transalation
//------------------------------------------------------------------------------

// Translate translates a gohttp.Request object into an http.Request object.
func (r *Request) Translate(client *Client) (*http.Request, error) {
	URL := client.BaseURL + r.URL

	var err error
	var req *http.Request
	if r.Body != nil {
		req, err = r.requestWithBody(r.Body, r.Method, URL)
		if err != nil {
			return nil, err
		}
	} else if r.Data != nil {
		req, err = r.requestWithData(r.Data, r.Method, URL)
		if err != nil {
			return nil, err
		}
	} else {
		req, err = http.NewRequest(r.Method, URL, nil)
		if err != nil {
			return nil, err
		}
	}

	// Hydrate go request with prams from gohttp.Request object.
	err = r.hydrateRequest(req, client)
	if err != nil {
		return nil, err
	}

	return req, nil
}

//------------------------------------------------------------------------------
// Helpers
//------------------------------------------------------------------------------

func (r *Request) hydrateRequest(req *http.Request, client *Client) error {
	// Add request headers.
	req.Header = r.headersForRequest(client.Headers)

	// Setup basic auth if needed.
	if client.BasicAuth != nil {
		req.SetBasicAuth(client.BasicAuth.Username, client.BasicAuth.Password)
	}

	// Add request parameters.
	if params := r.paramsForRequest(); params != "" {
		req.URL.RawQuery = params
	}
	return nil
}

func (r *Request) requestWithBody(body interface{}, method string, url string) (*http.Request, error) {
	jsonData, err := JSONData(body)
	if err != nil {
		return nil, err
	}

	return http.NewRequest(method, url, jsonData)
}

func (r *Request) requestWithData(data interface{}, method string, url string) (*http.Request, error) {
	formData, err := FormData(data)
	if err != nil {
		return nil, err
	}

	return http.NewRequest(method, url, formData)
}

func (r *Request) paramsForRequest() string {
	values := url.Values{}
	for key, value := range r.Params {
		values.Add(key, value)
	}
	return values.Encode()
}

func (r *Request) headersForRequest(headers http.Header) http.Header {
	for key, values := range r.Header {
		if len(headers.Get(key)) == 0 {
			for _, headerValue := range values {
				headers.Add(key, headerValue)
			}
		}
	}
	return headers
}
