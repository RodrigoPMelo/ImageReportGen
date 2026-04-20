package service

import (
	"errors"
	"path/filepath"
	"strings"

	"ImageReportGen/core/domain"
	"ImageReportGen/core/ports"
)

type ReportService struct {
	extractor     ports.FileExtractor
	imageProc     ports.ImageProcessor
	reportGen     ports.ReportGenerator
	tempOutputDir string
}

func NewReportService(
	extractor ports.FileExtractor,
	imageProc ports.ImageProcessor,
	reportGen ports.ReportGenerator,
	tempOutputDir string,
) *ReportService {
	return &ReportService{
		extractor:     extractor,
		imageProc:     imageProc,
		reportGen:     reportGen,
		tempOutputDir: tempOutputDir,
	}
}

func (s *ReportService) ProcessInputPaths(paths []string) (domain.ProcessedInput, error) {
	result := domain.ProcessedInput{
		ImagePaths:   make([]string, 0, len(paths)),
		IgnoredPaths: []string{},
	}

	for _, path := range paths {
		lower := strings.ToLower(path)
		switch {
		case strings.HasSuffix(lower, ".zip"):
			extracted, err := s.extractor.ExtractZip(path, s.tempOutputDir)
			if err != nil {
				return domain.ProcessedInput{}, err
			}
			result.ImagePaths = append(result.ImagePaths, extracted...)
		case isSupportedImage(lower):
			result.ImagePaths = append(result.ImagePaths, path)
		default:
			result.IgnoredPaths = append(result.IgnoredPaths, path)
		}
	}

	return result, nil
}

func (s *ReportService) GenerateReport(request domain.ReportRequest) (domain.ReportResult, error) {
	if len(request.ImagePaths) == 0 {
		return domain.ReportResult{}, errors.New("por favor, adicione pelo menos uma imagem")
	}

	outputPath := request.OutputPath
	if strings.TrimSpace(outputPath) == "" {
		outputPath = "relatorio_gerado.docx"
	}

	landscapeImages := []string{}
	portraitImages := []string{}
	for _, imagePath := range request.ImagePaths {
		orientation, err := s.imageProc.GetImageOrientation(imagePath)
		if err != nil {
			return domain.ReportResult{}, err
		}

		if orientation == domain.OrientationLandscape {
			landscapeImages = append(landscapeImages, imagePath)
			continue
		}
		portraitImages = append(portraitImages, imagePath)
	}

	err := s.reportGen.GenerateReport(request.TemplatePath, landscapeImages, portraitImages, outputPath)
	if err != nil {
		return domain.ReportResult{}, err
	}

	if s.tempOutputDir != "" {
		if cleanErr := s.extractor.CleanUpTempFiles(s.tempOutputDir); cleanErr != nil {
			return domain.ReportResult{}, cleanErr
		}
	}

	return domain.ReportResult{
		OutputPath:     filepath.Clean(outputPath),
		TotalImages:    len(request.ImagePaths),
		LandscapeCount: len(landscapeImages),
		PortraitCount:  len(portraitImages),
	}, nil
}

func isSupportedImage(path string) bool {
	return strings.HasSuffix(path, ".jpg") ||
		strings.HasSuffix(path, ".jpeg") ||
		strings.HasSuffix(path, ".png")
}
