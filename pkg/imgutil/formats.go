package imgutil

import (
	"fmt"

	"github.com/h2non/bimg"
)

// GetImageType converts string format to bimg type
func GetImageType(format string) (bimg.ImageType, error) {
	switch format {
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
	default:
		return bimg.JPEG, fmt.Errorf("unsupported format: %s", format)
	}
}
