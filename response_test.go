package room

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

// TestNewResponse tests the NewResponse function.
func TestNewResponse(t *testing.T) {
	// Mock HTTP response
	httpResp := &http.Response{
		StatusCode: http.StatusOK,
		Request: &http.Request{
			Method: http.MethodGet,
			URL: &url.URL{
				Scheme: "https",
				Host:   "example.com",
				Path:   "/test",
			},
		},
		Header: http.Header{},
		Body:   io.NopCloser(strings.NewReader("test body")),
	}

	// Test case for successful creation of Response
	resp, err := NewResponse(httpResp, false)
	if err != nil {
		t.Errorf("NewResponse() returned unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("NewResponse() returned response with unexpected status code: %d", resp.StatusCode)
	}

	// Test case for DTO creation when forceDTO is true
	resp, err = NewResponse(httpResp, true)
	if err != nil {
		t.Errorf("NewResponse(forceDTO=true) returned unexpected error: %v", err)
	}
	if resp.DTO == nil {
		t.Error("NewResponse(forceDTO=true) did not create DTO")
	}

	// Test case for error in NewResponse
	resp, err = NewResponse(httpResp, false)
	if err == nil {
		t.Error("NewResponse() did not return expected error for HTTP client error")
	}
}

// TestResponse_OK tests the OK method of the Response struct.
func TestResponse_OK(t *testing.T) {
	// Test case for successful response (status code 200)
	response := Response{StatusCode: 200}
	if !response.OK() {
		t.Error("Response OK() returned false for status code 200")
	}

	// Test case for unsuccessful response (status code 404)
	response = Response{StatusCode: 404}
	if response.OK() {
		t.Error("Response OK() returned true for status code 404")
	}
}

// TestResponse_SetHeader tests the SetHeader method of the Response struct.
func TestResponse_SetHeader(t *testing.T) {
	// Test case for setting header
	header := http.Header{}
	header.Add("Content-Type", "application/json")
	response := Response{}
	response = response.setHeader(header)
	if response.Header == nil {
		t.Error("Response SetHeader() did not set the header correctly")
	}
}

// TestResponse_SetRequestHeader tests the SetRequestHeader method of the Response struct.
func TestResponse_SetRequestHeader(t *testing.T) {
	// Test case for setting request header
	header := http.Header{}
	header.Add("Content-Type", "application/json")
	response := Response{}
	response = response.setRequestHeader(header)
	if response.RequestHeader == nil {
		t.Error("Response SetRequestHeader() did not set the request header correctly")
	}
}

// TestResponse_SetRequestBodyData tests the SetRequestBodyData method of the Response struct.
func TestResponse_SetRequestBodyData(t *testing.T) {
	// Test case for setting request body data
	reqBody := &strings.Reader{}
	httpReq := &http.Request{
		Body: io.NopCloser(reqBody),
	}
	response := Response{}
	response = response.setRequestBodyData(httpReq)
	if response.RequestBody == nil {
		t.Error("Response SetRequestBodyData() did not set the request body data correctly")
	}
}

// TestResponse_SetRequestURI tests the SetRequestURI method of the Response struct.
func TestResponse_SetRequestURI(t *testing.T) {
	// Test case for setting request URI
	httpReq := &http.Request{
		URL: &url.URL{
			Scheme: "https",
			Host:   "example.com",
			Path:   "/test",
		},
	}
	response := Response{}
	response = response.setRequestURI(httpReq)

	if response.RequestURI.Scheme() != "https" {
		t.Error("Response SetRequestURI() did not set the request URI correctly")
	}
}

// TestResponse_SetData tests the SetData method of the Response struct.
func TestResponse_SetData(t *testing.T) {
	// Test case for setting response data
	httpResp := &http.Response{
		Body: io.NopCloser(strings.NewReader("test data")),
	}
	response := Response{}
	response, err := response.setData(httpResp)
	if err != nil {
		t.Errorf("Response SetData() returned unexpected error: %v", err)
	}
	if len(response.Data) == 0 {
		t.Error("Response SetData() did not set the response data correctly")
	}
}

type TestDTO struct {
	Key string `json:"key"`
}

// TestResponse_SetDTO tests the SetDTO method of the Response struct.
func TestResponse_SetDTO(t *testing.T) {
	// Test case for setting DTO
	response := Response{Header: NewHeader(), Data: []byte(`{"key": "value"}`), DTO: &TestDTO{}}
	response, err := response.setDTO(false)
	if err != nil {
		t.Errorf("Response SetDTO() returned unexpected error: %v", err)
	}

	if response.DTO == nil {
		t.Error("Response SetDTO() did not set the DTO correctly")
	}
}
