package services

import (
	"github.com/h2non/bimg"

	"github.com/michaelyang12/frame/pkg/imgutil"
)

type ImageProcessor struct{}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{}
}

func (p *ImageProcessor) Resize(imageBytes []byte, width, height int, format string) ([]byte, string, error) {
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
		return nil, "", err
	}

	// Determine content type
	contentType := "image/jpeg" // Default
	if options.Type == bimg.PNG {
		contentType = "image/png"
	} else if options.Type == bimg.WEBP {
		contentType = "image/webp"
	} else if options.Type == bimg.TIFF {
		contentType = "image/tiff"
	}

	return newImage, contentType, nil
}

// func (p *ImageProcessor) Convert(imageBytes []byte, format string, quality int) ([]byte, string, models.ImageInfo, error) {
// 	// Implementation for conversion
// 	// ...
// }
