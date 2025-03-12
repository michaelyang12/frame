package services

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/h2non/bimg"

	"github.com/michaelyang12/frame/internal/models"
	"github.com/michaelyang12/frame/pkg/imgutil"
)

type ImageProcessor struct{}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{}
}

func (p *ImageProcessor) Resize(imageBytes []byte, width, height int, format string) ([]byte, string, models.ImageInfo, error) {
	options := bimg.Options{
		Width:  width,
		Height: height,
	}

	// Apply format conversion if specified
	if format != "" {
		imageType, err := imgutil.GetImageType(format)
		if err == nil {
			options.Type = imageType
		}
	}

	// Resize the image
	newImage, err := bimg.NewImage(imageBytes).Process(options)
	if err != nil {
		return nil, "", models.ImageInfo{}, err
	}

	metadata := getImageMetadata(newImage, format)

	// Determine content type
	contentType := imgutil.GetContentType(options.Type)
	return newImage, contentType, metadata, nil
}

func (p *ImageProcessor) Convert(imageBuffer *bytes.Buffer, inputContentType string, outputFormat string, quality int) ([]byte, string, models.ImageInfo, error) {
	// imageBytes := buffer.Bytes()
	if inputContentType == "image/heic" || inputContentType == "image/heif" {
		log.Printf("HEIC file detected \n")
	}

	imageType, err := imgutil.GetImageType(outputFormat)
	if err != nil {
		return nil, "", models.ImageInfo{}, err
	}

	options := bimg.Options{
		Type:    imageType,
		Quality: quality,
	}

	newImage, err := bimg.NewImage(imageBuffer.Bytes()).Process(options)
	if err != nil {
		return nil, "", models.ImageInfo{}, fmt.Errorf("Error processing image: %w", err)
	}
	metadata := getImageMetadata(newImage, outputFormat)
	// Set appropriate content type
	contentType := imgutil.GetContentType(options.Type)
	return newImage, contentType, metadata, nil
}

func getImageMetadata(img []byte, outputFormat string) models.ImageInfo {
	info, _ := bimg.NewImage(img).Size()
	return models.ImageInfo{
		Width:    info.Width,
		Height:   info.Height,
		Size:     len(img),
		Format:   outputFormat,
		Modified: time.Now().Format(time.RFC3339),
	}
}

// func ConvertHEICToJPEG(heicBuffer *bytes.Buffer) (*bytes.Buffer, error) {
// 	// Decode HEIC
// 	img, err := heif.Decode(heicBuffer)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to decode HEIC: %w", err)
// 	}

// 	// Convert to JPEG
// 	jpegBuffer := new(bytes.Buffer)
// 	err = jpeg.Encode(jpegBuffer, img, nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to encode JPEG: %w", err)
// 	}

// 	return jpegBuffer, nil
// }
