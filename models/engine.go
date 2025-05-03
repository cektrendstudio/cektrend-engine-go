package models

import "time"

type WebScreenshotRequest struct {
	Key string `json:"key"`
	URL string `json:"url"`
}

type WebScreenshotResponse struct {
	URL string `json:"url"`
}

type CreatePhishingWebReportRequest struct {
	SiteURL  string `json:"site_url"`
	ImageURL string `json:"image_url"`
	Version  int8   `json:"version"`
}

type GetPhisingWebReportResponse struct {
	ID        int64     `json:"id" db:"id"`
	SiteURL   string    `json:"site_url" db:"site_url"`
	ImageURL  string    `json:"image_url" db:"image_url"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
