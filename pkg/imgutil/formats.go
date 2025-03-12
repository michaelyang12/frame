package imgutil

import (
	"fmt"
	"strings"

	"github.com/h2non/bimg"
)

// GetImageType converts string format to bimg type
func GetImageType(format string) (bimg.ImageType, error) {
	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		return bimg.JPEG, nil
	case "png":
		return bimg.PNG, nil
	case "webp":
		return bimg.WEBP, nil
	case "tiff":
		return bimg.TIFF, nil
	case "gif":
		return bimg.GIF, nil
	case "heif", "heic":
		return bimg.HEIF, nil
	default:
		return bimg.JPEG, fmt.Errorf("unsupported format: %s", format)
	}
}

func GetContentType(imageType bimg.ImageType) string {
	contentType := "image/jpeg" // Default
	if imageType == bimg.PNG {
		contentType = "image/png"
	} else if imageType == bimg.WEBP {
		contentType = "image/webp"
	} else if imageType == bimg.TIFF {
		contentType = "image/tiff"
	} else if imageType == bimg.GIF {
		contentType = "image/gif"
	} else if imageType == bimg.HEIF {
		contentType = "image/heif"
	}
	return contentType
}
