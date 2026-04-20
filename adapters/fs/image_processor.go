package fs

import (
	"image"
	"os"

	"ImageReportGen/core/domain"
)

type ImageProcessor struct{}

func NewImageProcessor() *ImageProcessor {
	return &ImageProcessor{}
}

func (p *ImageProcessor) GetImageOrientation(imagePath string) (domain.Orientation, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	imgConfig, _, err := image.DecodeConfig(file)
	if err != nil {
		return "", err
	}

	if imgConfig.Width > imgConfig.Height {
		return domain.OrientationLandscape, nil
	}

	return domain.OrientationPortrait, nil
}
