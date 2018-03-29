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
	Params []Param

	// Body contains the JSON body to be used for the request. Body will be sent as application/json.
	Body interface{}

	// Form contains the form to be used for the request. Form will be sent as application/x-www-form-urlencoded.
	Form interface{}
}

// Param holds the key/value pair associated with a parameter on a Request
type Param struct {
	Key string

	Value string
}

//------------------------------------------------------------------------------
// Request Translation
//------------------------------------------------------------------------------

// Translate translates a `gohttp.Request` object into an `http.Request` object.
func (r *Request) Translate(client *Client) (*http.Request, error) {
	URL := client.BaseURL + r.URL

	var err error
	var req *http.Request
	if r.Body != nil {
		// Request with JSON body.
		req, err = r.requestWithBody(r.Body, r.Method, URL)
		if err != nil {
			return nil, err
		}
	} else if r.Form != nil {
		// Request with form data.
		req, err = r.requestWithForm(r.Form, r.Method, URL)
		if err != nil {
			return nil, err
		}
	} else {
		// Request without data.
		req, err = http.NewRequest(r.Method, URL, nil)
		if err != nil {
			return nil, err
		}
	}

	// Hydrate http.Request with details from gohttp.Request object.
	r.hydrateRequest(req, client)

	return req, nil
}

//------------------------------------------------------------------------------
// Params
//------------------------------------------------------------------------------

// SetParam adds a new request parameter.
func (r *Request) SetParam(key string, value string) {
	param := Param{
		Key:   key,
		Value: value,
	}
	r.Params = append(r.Params, param)
}

//------------------------------------------------------------------------------
// Helpers
//------------------------------------------------------------------------------

func (r *Request) requestWithBody(body interface{}, method string, url string) (*http.Request, error) {
	jsonData, err := JSONData(body)
	if err != nil {
		return nil, err
	}

	return http.NewRequest(method, url, jsonData)
}

func (r *Request) requestWithForm(form interface{}, method string, url string) (*http.Request, error) {
	formData, err := FormData(form)
	if err != nil {
		return nil, err
	}

	return http.NewRequest(method, url, formData)
}

func (r *Request) hydrateRequest(req *http.Request, client *Client) {
	// Set basic auth if needed.
	if client.BasicAuth != nil {
		req.SetBasicAuth(client.BasicAuth.Username, client.BasicAuth.Password)
	}

	// Add request parameters.
	if params := r.paramsForRequest(); params != "" {
		req.URL.RawQuery = params
	}

	// Add request headers.
	req.Header = r.combineClientHeaders(client.Headers)
}

func (r *Request) paramsForRequest() string {
	values := url.Values{}
	for _, param := range r.Params {
		values.Add(param.Key, param.Value)
	}
	return values.Encode()
}

func (r *Request) combineClientHeaders(headers http.Header) http.Header {
	if r.Header == nil {
		r.Header = headers
	} else {
		for key, values := range headers {
			if len(r.Header.Get(key)) == 0 {
				for _, headerValue := range values {
					r.Header.Add(key, headerValue)
				}
			}
		}
	}
	return r.Header
}
