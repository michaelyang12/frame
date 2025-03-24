package handlers

import (
	"net/http"
	"strconv"

	"github.com/michaelyang12/frame/internal/models"
	"github.com/michaelyang12/frame/internal/services"
)

// HandleResize processes image resize requests
func HandleResize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data and get file
	buffer, _, outputFormat, err := services.GetImageBufferFromFormData(w, r)
	if err != nil {
		http.Error(w, "Failed to get image from form data: "+err.Error(), http.StatusInternalServerError)
	}

	// Extract dimensions
	width, _ := strconv.Atoi(r.FormValue("width"))
	height, _ := strconv.Atoi(r.FormValue("height"))

	// If no form values, try URL params
	if width == 0 {
		width, _ = strconv.Atoi(r.URL.Query().Get("width"))
	}
	if height == 0 {
		height, _ = strconv.Atoi(r.URL.Query().Get("height"))
	}

	// Default dimensions
	if width == 0 {
		width = 800
	}
	if height == 0 {
		height = 600
	}

	// Process the image
	processor := services.NewImageProcessor()
	newImage, contentType, metadata, err := processor.Resize(buffer.Bytes(), width, height, outputFormat)
	if err != nil {
		http.Error(w, "Failed to process image: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Decide if we should return image or metadata based on query param
	returnMetadata := r.URL.Query().Get("metadata") == "true" || r.FormValue("metadata") == "true"

	if returnMetadata {
		resp := models.Response{
			Success: true,
			Message: "Image resized successfully",
			Data:    metadata,
		}
		sendJSONResponse(w, resp, http.StatusOK)
		return
	}

	// Set headers and write response
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(newImage)))
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Write(newImage)
}

// HandleConvert converts images between formats
func HandleConvert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data and get file
	buffer, inputContentType, outputFormat, err := services.GetImageBufferFromFormData(w, r)
	if err != nil {
		http.Error(w, "Failed to get image from form data: "+err.Error(), http.StatusInternalServerError)
	}

	// Get quality setting
	quality, _ := strconv.Atoi(r.FormValue("quality"))
	if quality == 0 {
		quality, _ = strconv.Atoi(r.URL.Query().Get("quality"))
	}
	if quality <= 0 || quality > 100 {
		quality = 80 // Default quality
	}

	// Process image
	processor := services.NewImageProcessor()
	newImage, contentType, metadata, err := processor.Convert(buffer, inputContentType, outputFormat, quality)
	if err != nil {
		http.Error(w, "Failed to convert image: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Decide if we should return image or metadata based on query param
	returnMetadata := r.URL.Query().Get("metadata") == "true" || r.FormValue("metadata") == "true"

	if returnMetadata {
		resp := models.Response{
			Success: true,
			Message: "Image converted successfully",
			Data:    metadata,
		}
		sendJSONResponse(w, resp, http.StatusOK)
		return
	}

	// Set headers and write response
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(newImage)))
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Write(newImage)
}

func HandleTrim(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data and get file
	buffer, _, outputFormat, err := services.GetImageBufferFromFormData(w, r)
	if err != nil {
		http.Error(w, "Failed to get image from form data: "+err.Error(), http.StatusInternalServerError)
	}

	// Process the image
	processor := services.NewImageProcessor()
	newImage, contentType, err := processor.RemoveBackground(buffer, outputFormat)
	if err != nil {
		http.Error(w, "Failed to remove background from image: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set headers and write response
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(newImage)))
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Write(newImage)
}
