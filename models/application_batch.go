package models

import "time"

type ApplicationBatch struct {
	BatchID               int64     `json:"batch_id" db:"batch_id"`
	BatchName             string    `json:"batch_name" db:"batch_name"`
	FileName              string    `json:"file_name" db:"file_name"`
	FileURL               string    `json:"file_url" db:"file_url"`
	ApplicationType       string    `json:"application_type" db:"application_type"`
	ProcessingStatus      string    `json:"processing_status" db:"processing_status"`
	ProcessingStartedAt   string    `json:"processing_started_at" db:"processing_started_at"`
	ProcessingCompletedAt string    `json:"processing_completed_at" db:"processing_completed_at"`
	TotalRow              int64     `json:"total_row" db:"total_row"`
	TotalSuccess          int64     `json:"total_success" db:"total_success"`
	Properties            string    `json:"properties" db:"properties"`
	StoreID               int64     `json:"store_id" db:"store_id"`
	UserID                int64     `json:"user_id" db:"user_id"`
	CreatedAt             time.Time `json:"created_at,omitempty" db:"created_at"`
	CreatedBy             string    `json:"created_by" db:"created_by"`
	UpdatedAt             time.Time `json:"updated_at,omitempty" db:"updated_at"`
	UpdatedBy             string    `json:"updated_by" db:"updated_by"`
}

type ApplicationBatchDetail struct {
	ID            int64     `json:"id" db:"id"`
	BatchID       int64     `json:"batch_id" db:"batch_id"`
	Order         int64     `json:"order" db:"order"`
	Status        string    `json:"status" db:"status"`
	FailedMessage string    `json:"failed_message" db:"failed_message"`
	Properties    string    `json:"properties" db:"properties"`
	StoreID       int64     `json:"store_id" db:"store_id"`
	UserID        int64     `json:"user_id" db:"user_id"`
	CreatedAt     time.Time `json:"created_at,omitempty" db:"created_at"`
	CreatedBy     string    `json:"created_by" db:"created_by"`
	UpdatedAt     time.Time `json:"updated_at,omitempty" db:"updated_at"`
	UpdatedBy     string    `json:"updated_by" db:"updated_by"`
}

type UploadApplicationRequest struct {
	UserID           int64  `json:"user_id"`
	StoreID          int64  `json:"store_id" validate:"required"`
	BatchName        string `json:"batch_name"`
	FileName         string `json:"file_name"`
	FileURL          string `json:"file_url"`
	Properties       string `json:"properties"`
	ApplicationType  string `json:"application_type"`
	ProcessingStatus string `json:"processing_status"`
	TotalRow         int64  `json:"total_row"`
	CreatedBy        string `json:"created_by"`
}
