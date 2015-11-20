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

	// Client is the underlying `http.Client` object that is used to issue
	// requests.
	Client *http.Client
}

func NewClient(baseURL string, headers http.Header) Client {
	return Client{baseURL, headers, &http.Client{}}
}

func (c *Client) Execute(r *Request) (map[string]interface{}, error) {
	switch r.Method {
	case GET:
		return c.Get(r.URL)
	case POST:
		return c.Post(r.URL, r.Body)
	case DELETE:
		return c.Delete(r.URL)
	case PATCH:
		return c.Patch(r.URL, r.Body)
	}
	return nil, nil
}

func (c *Client) Get(url string) (map[string]interface{}, error) {
	URL := c.BaseURL + url
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}

	req.Header = c.Headers
	return performRequest(req, c.Client)
}

func (c *Client) Post(url string, params interface{}) (map[string]interface{}, error) {
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
	return performRequest(req, c.Client)
}

func (c *Client) Delete(url string) (map[string]interface{}, error) {
	URL := c.BaseURL + url
	req, err := http.NewRequest("Delete", URL, nil)
	if err != nil {
		return nil, err
	}

	req.Header = c.Headers
	return performRequest(req, c.Client)
}

func (c *Client) Patch(url string, params map[string]interface{}) (map[string]interface{}, error) {
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
	return performRequest(req, c.Client)
}

func performRequest(r *http.Request, c *http.Client) (map[string]interface{}, error) {
	fmt.Printf("Performing request: %v\n", r)
	resp, err := c.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	fmt.Printf("Got response: %v\n", resp)
	if (resp.StatusCode != http.StatusOK) && (resp.StatusCode != http.StatusCreated) {
		return nil, fmt.Errorf("Unprocessable response code encountered: %v", resp.StatusCode)
	}
	return ParseJSON(resp.Body)
}

//------------------------------------------------------------------------------
// @Block Based Requests
//------------------------------------------------------------------------------

func (c *Client) PostWithBlock(r Request, s func(http.Request, map[string]interface{}), f func(http.Request, error)) {

}
