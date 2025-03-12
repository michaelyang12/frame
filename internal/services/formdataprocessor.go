package services

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/michaelyang12/frame/pkg/imgutil"
)

func GetImageBufferFromFormData(w http.ResponseWriter, r *http.Request) (*bytes.Buffer, string, string, error) {
	// Parse form data and get file
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		return nil, "", "", fmt.Errorf("Failed to parse form data: %w", err)
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		return nil, "", "", fmt.Errorf("Failed to read image: %w", err)
	}
	defer file.Close()

	sniffBuf := make([]byte, 512)
	if _, err := file.Read(sniffBuf); err != nil && err != io.EOF {
		return nil, "", "", fmt.Errorf("Failed to read file header: %w", err)
	}

	// Detect MIME type
	contentType := http.DetectContentType(sniffBuf)
	log.Printf("Detected MIME type: %s", contentType)

	//If cannot get MIME type, use file extensions
	if contentType == "application/octet-stream" {
		ext := filepath.Ext(header.Filename)
		filetype := strings.TrimPrefix(ext, ".")
		imageType, err := imgutil.GetImageType(filetype)
		if err != nil {
			return nil, "", "", fmt.Errorf("Failed to get image type: %w", err)
		}
		contentType = imgutil.GetContentType(imageType)
	}

	// Reset file cursor since we read first 512 bytes
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, "", "", fmt.Errorf("Failed to reset file: %w", err)
	}

	log.Printf("Converting file: %s, size: %d bytes, type: %s", header.Filename, header.Size, contentType)

	// Read file into buffer
	buffer := new(bytes.Buffer)
	if _, err := io.Copy(buffer, file); err != nil {
		return nil, "", "", fmt.Errorf("Failed to read file: %w", err)
	}

	// Extract format to convert to
	outputFormat := r.FormValue("format")
	if outputFormat == "" {
		outputFormat = r.URL.Query().Get("format")
	}

	return buffer, contentType, outputFormat, nil
}
