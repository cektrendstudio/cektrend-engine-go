package models

type WebScreenshotRequest struct {
	Key string `json:"key"`
	URL string `json:"url"`
}

type WebScreenshotResponse struct {
	URL string `json:"url"`
}
