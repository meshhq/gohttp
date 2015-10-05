package blaster

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Stereo struct {
	BaseURL string
	Headers http.Header
	Client  *http.Client
}

func NewStereo(baseURL string, headers http.Header) Stereo {
	return Stereo{baseURL, headers, &http.Client{}}
}

func (s *Stereo) Get(url string, params map[string]interface{}) (map[string]interface{}, error) {
	URL := s.BaseURL + url
	var req *http.Request
	var err error

	// If we don't have any params, we make request with no body
	if len(params) == 0 {
		req, err = http.NewRequest("GET", URL, nil)
	} else {
		jsonData, err := jsonData(params)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest("GET", URL, jsonData)
	}

	if err != nil {
		return nil, err

	}

	req.Header = s.Headers
	return performRequest(req, s.Client)
}

func (s *Stereo) Post(url string, params map[string]interface{}) (map[string]interface{}, error) {
	jsonData, err := jsonData(params)
	if err != nil {
		return nil, err
	}

	URL := s.BaseURL + url
	req, err := http.NewRequest("POST", URL, jsonData)
	if err != nil {
		return nil, err
	}

	req.Header = s.Headers
	return performRequest(req, s.Client)
}

func (s *Stereo) Delete(url string) (map[string]interface{}, error) {
	URL := s.BaseURL + url
	req, err := http.NewRequest("POST", URL, nil)
	req.Header = s.Headers
	resp, err := performRequest(req, s.Client)
	if err != nil {

	}
	return resp, err
}

func (s *Stereo) Patch(url string, params map[string]interface{}) (map[string]interface{}, error) {
	jsonData, err := json.Marshal(params)
	if err != nil {
		fmt.Printf("Failed to serialize json with error %v\n", err)
	}

	URL := s.BaseURL + url
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonData))
	req.Header = s.Headers
	resp, err := performRequest(req, s.Client)
	if err != nil {

	}
	return resp, err
}

func performRequest(r *http.Request, c *http.Client) (map[string]interface{}, error) {
	fmt.Printf("Request: %v\n", r)
	resp, err := c.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	fmt.Printf("Response: %v\n", resp)
	return ParseJSON(resp.Body)
}

func jsonData(params map[string]interface{}) (*bytes.Buffer, error) {
	jsonData, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(jsonData)
	return buffer, nil
}
