package gohttp

import (
	"fmt"
	"net/http"
)

const (
	// GET is a constant for the HTTP GET method.
	GET = "GET"

	// POST is a constant for the HTTP POST method.
	POST = "POST"

	// DELETE is a constant for the HTTP DELETE method.
	DELETE = "DELETE"

	// PUT is a constant for the HTTP PUT method.
	PUT = "PUT"

	// PATCH is a constant for the HTTP PATCH method.
	PATCH = "PATCH"
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
		response, err = c.Get(req.URL)
	case POST:
		response, err = c.Post(req.URL, req.Body)
	case DELETE:
		response, err = c.Delete(req.URL)
	case PUT:
		response, err = c.Put(req.URL, req.Body)
	case PATCH:
		response, err = c.Patch(req.URL, req.Body)
	}
	if response != nil {
		response.Request = req
	}
	return response, err
}

// Get performs an HTTP GET request with the supplied URL string.
func (c *Client) Get(url string) (*Response, error) {
	URL := c.BaseURL + url
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}

	req.Header = c.Headers
	return c.performRequest(req, c.Client)
}

// Post performs an HTTP POST request with the supplied URL string and
// parameters.
func (c *Client) Post(url string, params interface{}) (*Response, error) {
	jsonData, err := JSONData(params)
	if err != nil {
		return nil, err
	}

	URL := c.BaseURL + url
	req, err := http.NewRequest("POST", URL, jsonData)
	if err != nil {
		return nil, err
	}

	req.Header = c.Headers
	return c.performRequest(req, c.Client)
}

// Delete performs an HTTP DELETE request with the supplied URL string.
func (c *Client) Delete(url string) (*Response, error) {
	URL := c.BaseURL + url
	req, err := http.NewRequest("DELETE", URL, nil)
	if err != nil {
		return nil, err
	}

	req.Header = c.Headers
	return c.performRequest(req, c.Client)
}

// Put performs an HTTP PUT request with the supplied URL string and
// parameters.
func (c *Client) Put(url string, params interface{}) (*Response, error) {
	jsonData, err := JSONData(params)
	if err != nil {
		return nil, err
	}

	URL := c.BaseURL + url
	req, err := http.NewRequest("PUT", URL, jsonData)
	if err != nil {
		return nil, err
	}

	req.Header = c.Headers
	return c.performRequest(req, c.Client)
}

// Patch performs an HTTP PATCH request with the supplied URL string and
// parameters.
func (c *Client) Patch(url string, params interface{}) (*Response, error) {
	jsonData, err := JSONData(params)
	if err != nil {
		return nil, err
	}

	URL := c.BaseURL + url
	req, err := http.NewRequest("PATCH", URL, jsonData)
	if err != nil {
		return nil, err
	}

	req.Header = c.Headers
	return c.performRequest(req, c.Client)
}

func (c *Client) performRequest(r *http.Request, client *http.Client) (*Response, error) {
	if (c.BasicAuthUserName != "") && (c.BasicAuthPassword != "") {
		r.SetBasicAuth(c.BasicAuthUserName, c.BasicAuthPassword)
	}

	resp, err := client.Do(r)
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("Error performing request response: %v\n", err)
		return nil, err
	}

	response, err := NewResponse(resp)
	if (resp.StatusCode != http.StatusOK) && (resp.StatusCode != http.StatusCreated) && (resp.StatusCode != http.StatusNoContent) {
		return response, fmt.Errorf("Unprocessable response code encountered: %v", resp.StatusCode)
	}
	return response, err
}

//------------------------------------------------------------------------------
// @Block Based Requests
//------------------------------------------------------------------------------

// func (c *Client) PostWithBlock(r Request, s func(http.Request, map[string]interface{}), f func(http.Request, error)) {
//
// }
