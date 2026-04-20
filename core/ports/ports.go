package ports

import "ImageReportGen/core/domain"

type FileExtractor interface {
	ExtractZip(zipPath, outputDir string) ([]string, error)
	CleanUpTempFiles(tempDir string) error
}

type ImageProcessor interface {
	GetImageOrientation(imagePath string) (domain.Orientation, error)
}

type ReportGenerator interface {
	GenerateReport(templatePath string, landscapeImages, portraitImages []string, outputPath string) error
}
