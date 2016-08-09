package gohttp

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/meshhq/gohttp/Godeps/_workspace/src/gopkg.in/check.v1"
)

type ResponseTest struct{}

var _ = check.Suite(&ResponseTest{})

func (r *ResponseTest) SetUpSuite(c *check.C) {
	server = httptest.NewServer(RouteRequest())
}

func (r *ResponseTest) TestBuildingResponseWithJSON(c *check.C) {
	url := fmt.Sprintf("%v/json", server.URL)
	req, _ := http.NewRequest(POST, url, nil)

	client := &http.Client{}
	res, _ := client.Do(req)
	response, err := NewResponse(res)
	c.Assert(err, check.IsNil)
	c.Assert(response.Body, check.DeepEquals, map[string]interface{}{"test": "test"})
	c.Assert(response.Code, check.Equals, http.StatusCreated)
}
