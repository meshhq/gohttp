package gohttp

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/meshhq/meshCore/lib/gohttp/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/meshhq/meshCore/lib/gohttp/Godeps/_workspace/src/gopkg.in/check.v1"
	"github.com/meshhq/meshCore/lib/meshRedis"
)

type ClientTest struct{}

var _ = check.Suite(&ClientTest{})

var server *httptest.Server

var requestCounter = 0

func (r *ClientTest) SetUpSuite(c *check.C) {
	err := meshRedis.SetupRedis()
	c.Assert(err, check.Equals, nil)

	server = httptest.NewServer(RouteRequest())
}

func (r *ClientTest) SetUpTest(c *check.C) {
	requestCounter = 0
}

func (r *ClientTest) TearDownSuite(c *check.C) {
	err := meshRedis.ClosePool()
	c.Assert(err, check.Equals, nil)
}

//------------------------------------------------------------------------------
// Client Configuration
//------------------------------------------------------------------------------

func (r *ClientTest) TestInit(c *check.C) {
	client := NewClient("", nil)
	c.Assert(client, check.NotNil)
	c.Assert(client.BaseURL, check.NotNil)
	c.Assert(client.Headers, check.NotNil)
}

func (r *ClientTest) TestInitWithHeaders(c *check.C) {
	headers := http.Header{}
	client := NewClient("", headers)
	c.Assert(client.Headers, check.DeepEquals, headers)
}

func (r *ClientTest) TestInitWithBaseURL(c *check.C) {
	client := NewClient(server.URL, nil)
	c.Assert(client.BaseURL, check.Equals, server.URL)
}

func (r *ClientTest) TestBasicAuthRequest(c *check.C) {
	username := "testname"
	password := "testpass"

	client := NewClient(server.URL, nil)
	client.SetBasicAuth(username, password)
	c.Assert(client.BasicAuth, check.NotNil)
	c.Assert(client.BasicAuth.Username, check.Equals, username)
	c.Assert(client.BasicAuth.Password, check.Equals, password)
}

//------------------------------------------------------------------------------
// GET Request
//------------------------------------------------------------------------------

func (r *ClientTest) TestGetRequest(c *check.C) {
	client := NewClient(server.URL, nil)
	response, err := client.Execute(&Request{
		Method: GET,
		URL:    "/test",
	})
	c.Assert(err, check.IsNil)
	c.Assert(response.Code, check.Equals, http.StatusOK)
}

//------------------------------------------------------------------------------
// POST Request
//------------------------------------------------------------------------------

func (r *ClientTest) TestPostRequest(c *check.C) {
	client := NewClient(server.URL, nil)
	response, err := client.Execute(&Request{
		Method: POST,
		URL:    "/test",
	})
	c.Assert(err, check.IsNil)
	c.Assert(response.Code, check.Equals, http.StatusCreated)
}

func (r *ClientTest) TestPostRequestWithBody(c *check.C) {
	client := NewClient(server.URL, nil)
	response, err := client.Execute(&Request{
		Method: POST,
		URL:    "/test",
		Body:   map[string]interface{}{"name": "test"},
	})
	c.Assert(err, check.IsNil)
	c.Assert(response.Code, check.Equals, http.StatusCreated)
}

func (r *ClientTest) TestPostRequestWithForm(c *check.C) {
	client := NewClient(server.URL, nil)
	response, err := client.Execute(&Request{
		Method: POST,
		URL:    "/test",
		Data:   map[string]interface{}{"name": "test"},
	})
	c.Assert(err, check.IsNil)
	c.Assert(response.Code, check.Equals, http.StatusCreated)
}

//------------------------------------------------------------------------------
// PUT Request
//------------------------------------------------------------------------------
func (r *ClientTest) TestPutRequest(c *check.C) {
	client := NewClient(server.URL, nil)
	response, err := client.Execute(&Request{
		Method: PUT,
		URL:    "/test",
	})
	c.Assert(err, check.IsNil)
	c.Assert(response.Code, check.Equals, http.StatusNoContent)
}

func (r *ClientTest) TestPutRequestWithBody(c *check.C) {
	client := NewClient(server.URL, nil)
	response, err := client.Execute(&Request{
		Method: PUT,
		URL:    "/test",
		Body:   map[string]interface{}{"name": "test"},
	})
	c.Assert(err, check.IsNil)
	c.Assert(response.Code, check.Equals, http.StatusNoContent)
}

func (r *ClientTest) TestPutRequestWithForm(c *check.C) {
	client := NewClient(server.URL, nil)
	response, err := client.Execute(&Request{
		Method: PUT,
		URL:    "/test",
		Data:   map[string]interface{}{"name": "test"},
	})
	c.Assert(err, check.IsNil)
	c.Assert(response.Code, check.Equals, http.StatusNoContent)
}

//------------------------------------------------------------------------------
// PATH Request
//------------------------------------------------------------------------------

func (r *ClientTest) TestPatchRequest(c *check.C) {
	client := NewClient(server.URL, nil)
	response, err := client.Execute(&Request{
		Method: PATCH,
		URL:    "/test",
	})
	c.Assert(err, check.IsNil)
	c.Assert(response.Code, check.Equals, http.StatusNoContent)
}

func (r *ClientTest) TestPatchRequestWithBody(c *check.C) {
	client := NewClient(server.URL, nil)
	response, err := client.Execute(&Request{
		Method: PATCH,
		URL:    "/test",
		Body:   map[string]interface{}{"name": "test"},
	})
	c.Assert(err, check.IsNil)
	c.Assert(response.Code, check.Equals, http.StatusNoContent)
}

func (r *ClientTest) TestPatchRequestWithForm(c *check.C) {
	client := NewClient(server.URL, nil)
	response, err := client.Execute(&Request{
		Method: PATCH,
		URL:    "/test",
		Data:   map[string]interface{}{"name": "test"},
	})
	c.Assert(err, check.IsNil)
	c.Assert(response.Code, check.Equals, http.StatusNoContent)
}

//------------------------------------------------------------------------------
// DELETE Request
//------------------------------------------------------------------------------

func (r *ClientTest) TestDeleteRequest(c *check.C) {
	client := NewClient(server.URL, nil)
	response, err := client.Execute(&Request{
		Method: DELETE,
		URL:    "/test",
	})
	c.Assert(err, check.IsNil)
	c.Assert(response.Code, check.Equals, http.StatusNoContent)
}

//------------------------------------------------------------------------------
// Retry
//------------------------------------------------------------------------------

func (r *ClientTest) TestRetryScenario(c *check.C) {
	client := NewClient(server.URL, nil)
	client.RetryableStatusCodes = []int{500}
	_, err := client.Execute(&Request{
		Method: GET,
		URL:    "/retry",
	})
	c.Assert(err, check.NotNil)
	c.Assert(requestCounter > 1, check.Equals, true)
}

//------------------------------------------------------------------------------
// Retry
//------------------------------------------------------------------------------

func (r *ClientTest) TestRateLimitingClient(c *check.C) {
	client := NewClient(server.URL, nil)
	client.SetRateLimiterInfo(&RateLimitInfo{
		Token:        fmt.Sprintf("%v", time.Now()),
		MaxRequests:  10,   // 10 requests
		TimeInterval: 1000, // per second
	})

	// Execute 100 Requests
	for i := 0; i < 100; i++ {
		go func() {
			response, err := client.Execute(&Request{
				Method: GET,
				URL:    "/limit",
			})
			c.Assert(err, check.IsNil)
			c.Assert(response.Code, check.Equals, http.StatusOK)
		}()
	}
	time.Sleep(1 * time.Second)
	c.Assert(requestCounter, check.Equals, 10)

	time.Sleep(1 * time.Second)
	c.Assert(requestCounter, check.Equals, 20)
}

//------------------------------------------------------------------------------
// Test Router
//------------------------------------------------------------------------------

func RouteRequest() http.Handler {
	mux := mux.NewRouter()
	mux.HandleFunc("/test", HandleTest)
	mux.HandleFunc("/retry", HandleRetry)
	mux.HandleFunc("/limit", HandleLimit)
	return mux
}

func HandleTest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case GET:
		w.WriteHeader(http.StatusOK)

	case POST:
		w.WriteHeader(http.StatusCreated)

	case PUT, PATCH, DELETE:
		w.WriteHeader(http.StatusNoContent)
	}
}

func HandleRetry(w http.ResponseWriter, r *http.Request) {
	requestCounter++
	w.WriteHeader(http.StatusInternalServerError)
}

func HandleLimit(w http.ResponseWriter, r *http.Request) {
	requestCounter++
	fmt.Printf("Woo %v\n", requestCounter)
	w.WriteHeader(http.StatusOK)
}
