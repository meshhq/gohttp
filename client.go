package gohttp

import (
	"fmt"
	"net/http"
)

type Client struct {
	BaseURL string
	Headers http.Header
	Client  *http.Client
}

// HTTPRequest models an HTTP request sumbitted to the controller
// for processing.
type Request struct {

	// HTTP Method for the request
	Method string

	// Request URL
	URL string

	// Params models the body of the request as a hashmap.
	Body map[string]interface{}
}

func NewClient(baseURL string, headers http.Header) Client {
	return Client{baseURL, headers, &http.Client{}}
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

//------------------------------------------------------------------------------
// @Block Based Requests
//------------------------------------------------------------------------------

func (c *Client) PostWithBlock(r Request, s func(http.Request, map[string]interface{}), f func(http.Request, error)) {

}

func performRequest(r *http.Request, c *http.Client) (map[string]interface{}, error) {
	fmt.Printf("Performing Request: %v\n", r)
	resp, err := c.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	fmt.Printf("Response: %v\n", resp)
	if resp.StatusCode != http.StatusOK || resp.StatusCode != http.StatusCreated {
		return nil, err
	}
	return ParseJSON(resp.Body)
}
