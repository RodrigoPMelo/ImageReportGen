package service

import (
	"errors"
	"reflect"
	"testing"

	"ImageReportGen/core/domain"
)

type extractorMock struct {
	extractedByZip map[string][]string
	extractCalls   []string
	cleanupCalls   []string
	extractErr     error
	cleanupErr     error
}

func (m *extractorMock) ExtractZip(zipPath, outputDir string) ([]string, error) {
	m.extractCalls = append(m.extractCalls, zipPath+"|"+outputDir)
	if m.extractErr != nil {
		return nil, m.extractErr
	}
	return m.extractedByZip[zipPath], nil
}

func (m *extractorMock) CleanUpTempFiles(tempDir string) error {
	m.cleanupCalls = append(m.cleanupCalls, tempDir)
	return m.cleanupErr
}

type imageProcessorMock struct {
	orientations map[string]domain.Orientation
	err          error
}

func (m *imageProcessorMock) GetImageOrientation(imagePath string) (domain.Orientation, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.orientations[imagePath], nil
}

type reportGeneratorMock struct {
	lastTemplate  string
	lastLandscape []string
	lastPortrait  []string
	lastOutput    string
	err           error
}

func (m *reportGeneratorMock) GenerateReport(templatePath string, landscapeImages, portraitImages []string, outputPath string) error {
	if m.err != nil {
		return m.err
	}
	m.lastTemplate = templatePath
	m.lastLandscape = append([]string{}, landscapeImages...)
	m.lastPortrait = append([]string{}, portraitImages...)
	m.lastOutput = outputPath
	return nil
}

func TestProcessInputPathsExpandsZipAndFiltersSupportedFiles(t *testing.T) {
	extractor := &extractorMock{
		extractedByZip: map[string][]string{
			"set.zip": {"a.png", "b.jpeg"},
		},
	}
	svc := NewReportService(extractor, &imageProcessorMock{}, &reportGeneratorMock{}, "./temp_images")

	result, err := svc.ProcessInputPaths([]string{"set.zip", "photo.jpg", "notes.txt"})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	expectedImages := []string{"a.png", "b.jpeg", "photo.jpg"}
	if !reflect.DeepEqual(result.ImagePaths, expectedImages) {
		t.Fatalf("expected images %v, got %v", expectedImages, result.ImagePaths)
	}

	expectedIgnored := []string{"notes.txt"}
	if !reflect.DeepEqual(result.IgnoredPaths, expectedIgnored) {
		t.Fatalf("expected ignored %v, got %v", expectedIgnored, result.IgnoredPaths)
	}

	if len(extractor.extractCalls) != 1 || extractor.extractCalls[0] != "set.zip|./temp_images" {
		t.Fatalf("unexpected extract calls: %v", extractor.extractCalls)
	}
}

func TestGenerateReportReturnsErrorWhenImageListEmpty(t *testing.T) {
	svc := NewReportService(&extractorMock{}, &imageProcessorMock{}, &reportGeneratorMock{}, "./temp_images")

	_, err := svc.GenerateReport(domain.ReportRequest{
		TemplatePath: "template.docx",
		ImagePaths:   []string{},
		OutputPath:   "result.docx",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGenerateReportClassifiesCallsGeneratorAndCleanup(t *testing.T) {
	extractor := &extractorMock{}
	processor := &imageProcessorMock{
		orientations: map[string]domain.Orientation{
			"landscape.jpg": domain.OrientationLandscape,
			"portrait.jpg":  domain.OrientationPortrait,
		},
	}
	generator := &reportGeneratorMock{}

	svc := NewReportService(extractor, processor, generator, "./temp_images")

	result, err := svc.GenerateReport(domain.ReportRequest{
		TemplatePath: "template.docx",
		ImagePaths:   []string{"landscape.jpg", "portrait.jpg"},
		OutputPath:   "relatorio_gerado.docx",
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if generator.lastTemplate != "template.docx" {
		t.Fatalf("expected template template.docx, got %s", generator.lastTemplate)
	}
	if !reflect.DeepEqual(generator.lastLandscape, []string{"landscape.jpg"}) {
		t.Fatalf("unexpected landscape images: %v", generator.lastLandscape)
	}
	if !reflect.DeepEqual(generator.lastPortrait, []string{"portrait.jpg"}) {
		t.Fatalf("unexpected portrait images: %v", generator.lastPortrait)
	}
	if generator.lastOutput != "relatorio_gerado.docx" {
		t.Fatalf("expected output relatorio_gerado.docx, got %s", generator.lastOutput)
	}

	if result.LandscapeCount != 1 || result.PortraitCount != 1 || result.TotalImages != 2 {
		t.Fatalf("unexpected result counts: %+v", result)
	}

	if len(extractor.cleanupCalls) != 1 || extractor.cleanupCalls[0] != "./temp_images" {
		t.Fatalf("unexpected cleanup calls: %v", extractor.cleanupCalls)
	}
}

func TestGenerateReportPropagatesDependenciesErrors(t *testing.T) {
	svc := NewReportService(
		&extractorMock{},
		&imageProcessorMock{err: errors.New("orientation fail")},
		&reportGeneratorMock{},
		"./temp_images",
	)

	_, err := svc.GenerateReport(domain.ReportRequest{
		TemplatePath: "template.docx",
		ImagePaths:   []string{"x.jpg"},
		OutputPath:   "relatorio_gerado.docx",
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
