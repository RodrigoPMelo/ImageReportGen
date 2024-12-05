package utils

import (
	"image"
	"os"
)

func getImageOrientation(imagePath string) (string, error) {
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
		return "Paisagem", nil
	} else {
		return "Retrato", nil
	}
}