package gohttp

import (
	"fmt"
	"net/http"

	"gopkg.in/check.v1"
)

type RequestTest struct{}

var _ = check.Suite(&RequestTest{})

func (r *RequestTest) TestTransalatingJSONRequest(c *check.C) {
	header := http.Header{}
	header.Add(ContentType, "application/json")

	client := NewClient("", header)
	url := "api.google.com"
	method := POST
	body := map[string]interface{}{"test": "body"}
	params := map[string]string{"test": "params"}
	request := Request{
		URL:    url,
		Method: method,
		Body:   body,
		Params: params,
	}

	translated, err := request.Translate(client)
	c.Assert(err, check.Equals, nil)
	c.Assert(translated.Method, check.Equals, method)
	c.Assert(translated.URL.Path, check.Equals, url)
	c.Assert(translated.URL.RawQuery, check.Equals, fmt.Sprintf("%v=%v", "test", "params"))
	c.Assert(translated.Header, check.DeepEquals, header)
}

func (r *RequestTest) TestTransalatingFormRequest(c *check.C) {
	header := http.Header{}
	header.Add(ContentType, "application/json")

	client := NewClient("", header)
	url := "api.google.com"
	method := POST
	form := map[string]interface{}{"test": "body"}
	params := map[string]string{"test": "params"}
	request := Request{
		URL:    url,
		Method: method,
		Form:   form,
		Params: params,
	}

	translated, err := request.Translate(client)
	c.Assert(err, check.Equals, nil)
	c.Assert(translated.Method, check.Equals, method)
	c.Assert(translated.URL.Path, check.Equals, url)
	c.Assert(translated.URL.RawQuery, check.Equals, fmt.Sprintf("%v=%v", "test", "params"))
	c.Assert(translated.Header, check.DeepEquals, header)
}

func (r *RequestTest) TestTransalatingRequestWithBasicAuth(c *check.C) {
	client := NewClient("", nil)
	client.SetBasicAuth("username", "password")

	url := "api.google.com"
	method := GET
	request := Request{
		URL:    url,
		Method: method,
	}

	translated, err := request.Translate(client)
	c.Assert(err, check.Equals, nil)
	c.Assert(translated.Method, check.Equals, method)
	c.Assert(translated.URL.Path, check.Equals, url)
	c.Assert(translated.Header.Get("Authorization"), check.NotNil)
}

func (r *RequestTest) TestRequestWithJSONBody(c *check.C) {
	url := "api.google.com"
	method := POST
	body := map[string]interface{}{"test": "body"}

	request := Request{}
	goRequest, err := request.requestWithBody(body, method, url)
	c.Assert(err, check.Equals, nil)
	c.Assert(goRequest.Method, check.Equals, method)
	c.Assert(goRequest.URL.Path, check.Equals, url)
}

func (r *RequestTest) TestRequestWithFormData(c *check.C) {
	url := "api.google.com"
	method := POST
	form := map[string]interface{}{"test": "body"}

	request := Request{}
	goRequest, err := request.requestWithForm(form, method, url)
	c.Assert(err, check.Equals, nil)
	c.Assert(goRequest.Method, check.Equals, method)
	c.Assert(goRequest.URL.Path, check.Equals, url)
}

func (r *RequestTest) TestCombiningHeaders(c *check.C) {
	header := http.Header{}
	header.Add(ContentType, "application/json")

	requestHeader := http.Header{}
	requestHeader.Add(Accept, "application/json")
	request := Request{
		Header: requestHeader,
	}

	combined := request.combineClientHeaders(header)
	c.Assert(combined.Get(ContentType), check.NotNil)
	c.Assert(combined.Get(Accept), check.NotNil)
}
