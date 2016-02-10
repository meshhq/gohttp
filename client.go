package gohttp

import (
	"fmt"
	"net/http"
)

const (
	// HTTP Methods
	GET    = "GET"
	POST   = "POST"
	DELETE = "DELETE"
	PUT    = "PUT"
	PATCH  = "PATCH"

	// HTTP Headers
	CONTENT_TYPE  = "Content-Type"
	ACCEPT        = "Accept"
	AUTHORIZATION = "Authorization"
)

// Client models an HTTP client
type Client struct {

	// BaseURL models the base URL to be used for all requests issues by
	// the Client.
	BaseURL string

	// Headers model the global headers to be used for all requests issued
	// by the client.
	Headers http.Header

	BasicAuthUserName string
	BasicAuthPassword string

	// Client is the underlying `http.Client` object that is used to issue
	// requests.
	Client *http.Client
}

// NewClient instantiates a new instance of a gohttp.Client.
func NewClient(baseURL string, headers http.Header) *Client {
	if headers == nil {
		headers = http.Header{}
	}
	client := new(Client)
	client.BaseURL = baseURL
	client.Headers = headers
	client.Client = &http.Client{}
	return client
}

func (c *Client) SetHeader(header string, value string) {
	c.Headers.Add(header, value)
}

func (c *Client) SetBasicAuth(username string, password string) {
	c.BasicAuthUserName = username
	c.BasicAuthPassword = password
}

// Execute executes the HTTP request described with the given `gohttp.Request`.
func (c *Client) Execute(req *Request) (*Response, error) {
	var response *Response
	var err error
	switch req.Method {
	case GET:
		response, err = c.Get(req)
	case POST:
		response, err = c.Post(req)
	case DELETE:
		response, err = c.Delete(req)
	case PUT:
		response, err = c.Put(req)
	case PATCH:
		response, err = c.Patch(req)
	}
	if response != nil {
		response.Request = req
	}
	return response, err
}

// Get performs an HTTP GET request with the supplied URL string.
func (c *Client) Get(request *Request) (*Response, error) {
	URL := c.BaseURL + request.URL
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}

	// Add any request parameters
	query := req.URL.Query()
	for key, value := range request.Params {
		query.Add(key, value)
	}
	req.URL.RawQuery = query.Encode()
	req.Header = c.Headers
	return c.performRequest(req, c.Client)
}

// Post performs an HTTP POST request with the supplied URL string and
// parameters.
func (c *Client) Post(request *Request) (*Response, error) {
	jsonData, err := JSONData(request.Body)
	if err != nil {
		return nil, err
	}

	URL := c.BaseURL + request.URL
	req, err := http.NewRequest("POST", URL, jsonData)
	if err != nil {
		return nil, err
	}

	// Add any request parameters
	query := req.URL.Query()
	for key, value := range request.Params {
		query.Add(key, value)
	}
	req.URL.RawQuery = query.Encode()
	req.Header = c.Headers
	return c.performRequest(req, c.Client)
}

// Delete performs an HTTP DELETE request with the supplied URL string.
func (c *Client) Delete(request *Request) (*Response, error) {
	URL := c.BaseURL + request.URL
	req, err := http.NewRequest("DELETE", URL, nil)
	if err != nil {
		return nil, err
	}

	// Add any request parameters
	query := req.URL.Query()
	for key, value := range request.Params {
		query.Add(key, value)
	}
	req.URL.RawQuery = query.Encode()
	req.Header = c.Headers
	return c.performRequest(req, c.Client)
}

// Put performs an HTTP PUT request with the supplied URL string and
// parameters.
func (c *Client) Put(request *Request) (*Response, error) {
	jsonData, err := JSONData(request.Body)
	if err != nil {
		return nil, err
	}

	URL := c.BaseURL + request.URL
	req, err := http.NewRequest("PUT", URL, jsonData)
	if err != nil {
		return nil, err
	}

	// Add any request parameters
	query := req.URL.Query()
	for key, value := range request.Params {
		query.Add(key, value)
	}
	req.URL.RawQuery = query.Encode()
	req.Header = c.Headers
	return c.performRequest(req, c.Client)
}

// Patch performs an HTTP PATCH request with the supplied URL string and
// parameters.
func (c *Client) Patch(request *Request) (*Response, error) {
	jsonData, err := JSONData(request.Body)
	if err != nil {
		return nil, err
	}

	URL := c.BaseURL + request.URL
	req, err := http.NewRequest("PATCH", URL, jsonData)
	if err != nil {
		return nil, err
	}

	// Add any request parameters
	query := req.URL.Query()
	for key, value := range request.Params {
		query.Add(key, value)
	}
	req.URL.RawQuery = query.Encode()
	req.Header = c.Headers
	return c.performRequest(req, c.Client)
}

func (c *Client) performRequest(r *http.Request, client *http.Client) (*Response, error) {
	if (c.BasicAuthUserName != "") && (c.BasicAuthPassword != "") {
		r.SetBasicAuth(c.BasicAuthUserName, c.BasicAuthPassword)
	}
	fmt.Printf("Request: %v\n", r)
	resp, err := client.Do(r)
	if err != nil {
		fmt.Printf("Error performing request response: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	response, err := NewResponse(resp)
	return response, err
}

//------------------------------------------------------------------------------
// @Block Based Requests
//------------------------------------------------------------------------------

// func (c *Client) PostWithBlock(r Request, s func(http.Request, map[string]interface{}), f func(http.Request, error)) {
//
// }
