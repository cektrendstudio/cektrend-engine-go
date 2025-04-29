package models

type WebScreenshotRequest struct {
	URL string `json:"url" validate:"required,url"`
}

type WebScreenshotResponse struct {
	URL string `json:"url"`
}
