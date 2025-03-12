package handlers

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/h2non/bimg"
	"github.com/michaelyang12/frame/internal/models"
	"github.com/michaelyang12/frame/internal/services"
	"github.com/michaelyang12/frame/pkg/imgutil"
)

// HandleResize processes image resize requests
func HandleResize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form data and get file
	if err := r.ParseMultipartForm(32 << 20); err != nil { // 32MB max memory
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to read image: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Printf("Processing file: %s, size: %d bytes", header.Filename, header.Size)

	// Read file into buffer
	buffer := new(bytes.Buffer)
	if _, err := io.Copy(buffer, file); err != nil {
		http.Error(w, "Failed to read file: "+err.Error(), http.StatusInternalServerError)
		return
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

	// Get image format option
	format := r.FormValue("format")
	if format == "" {
		format = r.URL.Query().Get("format")
	}

	// Process the image
	processor := services.NewImageProcessor()
	newImage, contentType, err := processor.Resize(buffer.Bytes(), width, height, format)
	if err != nil {
		http.Error(w, "Failed to process image: "+err.Error(), http.StatusInternalServerError)
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
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to read image: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Printf("Converting file: %s, size: %d bytes", header.Filename, header.Size)

	// Read file into buffer
	buffer := new(bytes.Buffer)
	if _, err := io.Copy(buffer, file); err != nil {
		http.Error(w, "Failed to read file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract format to convert to
	format := r.FormValue("format")
	if format == "" {
		format = r.URL.Query().Get("format")
	}

	// Default to JPEG if no format specified
	if format == "" {
		format = "jpeg"
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
	imageBytes := buffer.Bytes()
	imageType, err := imgutil.GetImageType(format)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	options := bimg.Options{
		Type:    imageType,
		Quality: quality,
	}

	newImage, err := bimg.NewImage(imageBytes).Process(options)
	if err != nil {
		http.Error(w, "Failed to convert image: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get info about processed image
	info, _ := bimg.NewImage(newImage).Size()
	metadata := models.ImageInfo{
		Width:    info.Width,
		Height:   info.Height,
		Size:     len(newImage),
		Format:   format,
		Modified: time.Now().Format(time.RFC3339),
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

	// Set appropriate content type
	contentType := "image/jpeg" // Default
	if options.Type == bimg.PNG {
		contentType = "image/png"
	} else if options.Type == bimg.WEBP {
		contentType = "image/webp"
	} else if options.Type == bimg.TIFF {
		contentType = "image/tiff"
	} else if options.Type == bimg.GIF {
		contentType = "image/gif"
	}

	// Set headers and write response
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(newImage)))
	w.Header().Set("Cache-Control", "public, max-age=86400")
	w.Write(newImage)
}
