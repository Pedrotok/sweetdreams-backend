package controller

import (
	"encoding/json"
	"net/http"
)

type (
	// Response is the http json response schema
	Response struct {
		Status  int         `json:"status"`
		Message string      `json:"message"`
		Content interface{} `json:"content"`
	}

	// PaginatedResponse is the paginated response json schema
	// we not use it yet
	PaginatedResponse struct {
		Count    int         `json:"count"`
		Next     string      `json:"next"`
		Previous string      `json:"previous"`
		Results  interface{} `json:"results"`
	}
)

// NewResponse is the Response struct factory function.
func NewResponse(status int, message string, content interface{}) *Response {
	return &Response{
		Status:  status,
		Message: message,
		Content: content,
	}
}

// NewPaginatedResponse will created http paginated response
func NewPaginatedResponse(status, count int, message, next, prev string, results interface{}) *Response {
	return &Response{
		Status:  status,
		Message: message,
		Content: &PaginatedResponse{
			Count:    count,
			Next:     next,
			Previous: prev,
			Results:  results,
		},
	}
}

// ResponseWriter will write result in http.ResponseWriter
func ResponseWriter(res http.ResponseWriter, statusCode int, message string, data interface{}) error {
	res.WriteHeader(statusCode)
	httpResponse := NewResponse(statusCode, message, data)
	err := json.NewEncoder(res).Encode(httpResponse)
	return err
}
