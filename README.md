![alt text](Assets/gohttp.png)
[![Build Status](https://travis-ci.org/meshhq/gohttp.svg?branch=master)](https://travis-ci.org/meshhq/gohttp)

An HTTP networking client written in go.

`GoHTTP` is built on top of the `net/http` package and is designed to make integrating with JSON APIs from `golang` applications simple and easy.

## Why another HTTP Client?

`GoHTTP` offer two distinct, out-of-the-box features that separate it from other networking clients.

1. **Rate Limiting** - Most public APIs are going to have some form of rate limiting requirements. GoHTTP leverages [Funnel](github.com/meshhq/funnel) to make conforming to these rate limits trivial.
2. **Retry w/ Exponential BackOff** - Systems can fail or error in the wild and applications should be resilient. `GoHTTP` makes it easy for applications to specify the HTTP status codes that should result in a retry. `GoHTTP` retries using an exponential backoff algorithm powered by [Exponential Backoff](https://github.com/cenk/backoff).

`GoHTTP` also provides syntactic sugar on top of the `net/http` package that makes networking easier from go applications.

## Features

- [x] HTTP verbs supported Get, Put, Post, Patch, Delete
- [x] Rich object models for requests and responses
- [x] JSON and Form Data requests
- [x] Header and method constants
- [x] Retry with exponential backoff
- [x] Rate limiting backed by Redis
- [x] Response JSON parsing
- [x] JSON pretty printing
- [x] Full documentation
- [x] Comprehensive Unit Test Coverage

## Install

```
$ go get github.com/meshhq/gohttp
```

## Import

```
import "github.com/meshhq/gohttp"
```

## Examples

GET request

```go
request := gohttp.Request{Method: gohttp.GET}
client := &gohttp.NewClient("api.google.com", nil)
response, err := client.Execute(request)
```

POST request

```go
request := &gohttp.Request{
	Method: gohttp.POST,
	URL:    "/users",
	Body: 	map[string]interface{}{"first_name": "Kevin"},
}
client := gohttp.NewClient("api.google.com", nil)
response, err := client.Execute(request)
```

## Documentation
---

### Client

All requests are executed with a `Client`. `Client` objects can hold global parameters such as a `baseURL` or a set of `headers`. Global parameters will be applied to each request that is executed by the `Client`.

```go
baseURL := "api.google.com"
headers := map[string]string{gohttp.Accept : "application/json"}
client := gohttp.NewClient(baseURL, headers)
```

### Request

`gohttp` provides a `Request` object which makes building HTTP request simple and readable. The `URL` parameter of a request is relative to the `baseURL` parameter of a the `Client` that executes it.

```go
request := &gohttp.Request{
	Method: gohttp.POST,
	URL:    "/users",
	Body: 	map[string]interface{}{"first_name": "Kevin"},
}
response, err := client.Execute(request)
```

### Request Execution

A call to the `Execute` method on a `Client` object will return a `Response` object and an `error` that describes any failure that occurred.

```go
response, err := client.Execute(request)
if err != nil {
        return err
}
fmt.Printf("Response: %v", response)
```

### Response  

The `Response` object contains parsed information about the outcome of an HTTP request.

```go
fmt.Printf("Code: %v\n", resonse.Code) 		// int containing the response code.
fmt.Printf("Body: %v\n", resonse.Body) 		// interface{} containing the parsed response body.
fmt.Printf("Request: %v\n", resonse.Request) 	// `gohttp.Request` object which is a pointer to the original request.
```

#### Pretty Printing

Applications can also pretty print response objects via the `GoHTTP` convenience method `PrettyPrint`.

```go
gohttp.PrettyPrint(response)
```

### Retry

`GoHTTP` implements sophisticated retry logic with exponential backoff. Applications can configure the status codes for which an application should retry a request via the `RetryableStatusCodes` parameter on a `Client` object.

```
client.RetryableStatusCodes = []int{403} // Rate Limit Exceeded
```

Underneath the hood, `GoHTTP` leverages the [Exponential Backoff](https://github.com/cenk/backoff) package for building and executing the backoff algorithm. `GoHTTP` supplies a default backoff algorithm implementation, but applications can supply their own via the `Backoff` parameter on a `Client` object.

```
var backOff := func() error {
	// Algorithm configuration
}
client.Backoff = backOff
```

### Rate Limiting

GoHTTP provides robust, out of the box rate limiting support. This makes it very easy to conform with rate limiting policies published by APIs. Underneath the hood, `GoHTTP` leverages [Funnel](https://github.com/meshhq/funnel), a distributed rate limiter backed by redis.

#### Limit Info

Applications can configure their rate limiting policy by supplying a `LimitInfo` object to a `Client` object.

* `token` - A unique token upon which requests are limited. For example, if all requests to the gmail API need to be limited, applications could use the token "gmail".
* `MaxRequests` - The maximum number of requests that can take place within a given time interval.
* `TimeInterval` - The time interval in (milliseconds) to which the `MaxRequests` parameter applies.   

```go
info := &RateLimitInfo{
	Token:        "unique-token",
	MaxRequests:  10,
	TimeInterval: 1000,
}
client.SetRateLimiterInfo(info)
```

Applications can concurrently execute requests to the same client from as many `goroutines` as they wish. The `Client` will handle queuing requests with redis, and ensure that the rate limit is not breached and all requests are executed.

### Contributing
PRs are welcome, but will be rejected unless test coverage is updated
- [Taylor Halliday](https://github.com/tayhalla)
- [Kevin Coleman](https://github.com/kcoleman731)
