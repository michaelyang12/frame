package models

// Response represents a standardized API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ImageInfo stores processed image information
type ImageInfo struct {
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Size     int    `json:"size"`
	Format   string `json:"format"`
	Modified string `json:"modified"`
}
