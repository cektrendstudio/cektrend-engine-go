package models

import (
	"github.com/aws/aws-sdk-go/aws/session"
)

type ResponseSuccess struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}

type ResponseError struct {
	Message string      `json:"message"`
	Error   interface{} `json:"errors,omitempty"`
}

type PaginationResponse struct {
	CurrentPage  int   `json:"current_page"`
	PageSize     int   `json:"page_size"`
	TotalCount   int64 `json:"total_count"`
	TotalPages   int   `json:"total_pages"`
	FirstPage    int   `json:"first_page"`
	NextPage     int   `json:"next_page"`
	LastPage     int   `json:"last_page"`
	CurrentCount int   `json:"current_count"`
}

type RequestFilter struct {
	Select    []string
	Where     map[string]interface{}
	WhereOr   map[string]interface{}
	WhereIn   map[string][]interface{}
	WhereOrIn map[string][]interface{}
	Order     string
	Limit     int
	Offset    int
}

type AWSConfig struct {
	S3Session *session.Session
	Endpoint  string
	Bucket    string
	Region    string
}
