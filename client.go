package gohttp

import (
	"errors"
	"net/http"

	"github.com/cenk/backoff"
	"github.com/meshhq/funnel"
)

// HTTP Methods
const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	PATCH  = "PATCH"
	DELETE = "DELETE"
)

// HTTP Header Constants
const (
	ContentType   = "Content-Type"
	ContentLength = "Content-Length"
	Accept        = "Accept"
	Authorization = "Authorization"
	UserAgent     = "User-Agent"
)

// Client models an HTTP client.
//
// A GoHTTP Glient can contain global request parameters, such as a BaseURL,
// a set of Headers, or Basic Authentication Credentials. All global parameters
// will be applied to each request executed by the client.
//
// The GoHTTP client also has built in support support for retry (with
// exponetial backoff and rate limiting.
type Client struct {

	// BaseURL is base URL used for all requests executed by the client.
	BaseURL string

	// Headers model the global headers to be used for all requests issued by the client.
	Headers http.Header

	// BasicAuth
	BasicAuth *BasicAuth

	// RetryableStatusCodes is an array of codes that are retryable.
	RetryableStatusCodes []int

	// Backoff ...
	Backoff *backoff.ExponentialBackOff

	// RateLimiter is a rate limiter.
	rateLimiter *funnel.RateLimiter

	// goClient is the underlying `http.Client` that is used to issue requests.
	goClient *http.Client
}

//------------------------------------------------------------------------------
// Initailization
//------------------------------------------------------------------------------

// NewClient instantiates a new instance of a gohttp.Client.
func NewClient(baseURL string, headers http.Header) *Client {
	if headers == nil {
		headers = http.Header{}
	}

	client := new(Client)
	client.BaseURL = baseURL
	client.Headers = headers
	client.goClient = &http.Client{}
	client.Backoff = Backoff()
	client.RetryableStatusCodes = []int{http.StatusRequestTimeout, 429, 500}
	return client
}

//------------------------------------------------------------------------------
// Basic Authentication
//------------------------------------------------------------------------------

// SetBasicAuth configures basic authentication credentials for the client.
func (c *Client) SetBasicAuth(username string, password string) {
	basicAuth := &BasicAuth{
		Username: username,
		Password: password,
	}
	c.BasicAuth = basicAuth
}

//------------------------------------------------------------------------------
// Rate Limiting
//------------------------------------------------------------------------------

// SetRateLimiterInfo provides..
func (c *Client) SetRateLimiterInfo(limitInfo *funnel.RateLimitInfo) error {
	limiter, err := funnel.NewLimiter(limitInfo)
	if err != nil {
		return err
	}
	c.rateLimiter = limiter
	return nil
}

//------------------------------------------------------------------------------
// Request Execution
//------------------------------------------------------------------------------

// Execute executes the HTTPc request described with the given `gohttp.Request`.
func (c *Client) Execute(req *Request) (*Response, error) {
	var err error
	var response *Response

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

//------------------------------------------------------------------------------
// Convenience
//------------------------------------------------------------------------------

// Get performs an HTTP GET request with the supplied request object.
func (c *Client) Get(request *Request) (*Response, error) {
	req, err := request.Translate(c)
	if err != nil {
		return nil, err
	}
	return c.executeRequest(req)
}

// Post performs an HTTP POST request with the supplied request object.
func (c *Client) Post(request *Request) (*Response, error) {
	req, err := request.Translate(c)
	if err != nil {
		return nil, err
	}
	return c.executeRequest(req)
}

// Delete performs an HTTP DELETE request with the supplied request object.
func (c *Client) Delete(request *Request) (*Response, error) {
	req, err := request.Translate(c)
	if err != nil {
		return nil, err
	}
	return c.executeRequest(req)
}

// Put performs an HTTP PUT request with the supplied URL string and
// parameters.
func (c *Client) Put(request *Request) (*Response, error) {
	req, err := request.Translate(c)
	if err != nil {
		return nil, err
	}
	return c.executeRequest(req)
}

// Patch performs an HTTP PATCH request with the supplied URL string and
// parameters.
func (c *Client) Patch(request *Request) (*Response, error) {
	req, err := request.Translate(c)
	if err != nil {
		return nil, err
	}
	return c.executeRequest(req)
}

//------------------------------------------------------------------------------
// Request Execution
//------------------------------------------------------------------------------

func (c *Client) executeRequest(req *http.Request) (*Response, error) {

	var parsedError error
	var parsedResponse *Response

	// Setup our retryable operation.
	retry := func() error {

		// Execute the actual request.
		response, err := c.goClient.Do(req)
		if err != nil {
			parsedError = err
			return nil
		}

		// Defer closing response body if it exists.
		if response != nil {
			defer response.Body.Close()
		}

		// Parse our response into a gohttp.Response object.
		parsedResponse, parsedError = NewResponse(response)
		if parsedError != nil {
			return nil
		}

		// Retry status code checks.
		//
		// If the client encounters a retryable status code, we return
		// an error to signal a retry should occur
		for _, code := range c.RetryableStatusCodes {
			if code == parsedResponse.Code {
				return errors.New("Encountered retryable status code.")
			}
		}
		return nil
	}

	// Rate Limiter - If we are using a rate limiter, we enter here.
	//
	// This will block until the rate limiter is satisfied.
	if c.rateLimiter != nil {
		err := c.rateLimiter.Enter()
		if err != nil {
			return nil, err
		}
	}

	// Execute the retryable operation
	err := backoff.Retry(retry, c.Backoff)
	if err != nil {
		return nil, err
	}

	return parsedResponse, parsedError
}
