package main

import (
	"errors"
	"testing"

	"ImageReportGen/core/domain"
)

type reportUsecaseMock struct {
	processResult domain.ProcessedInput
	processErr    error

	lastGenerateRequest domain.ReportRequest
	generateResult      domain.ReportResult
	generateErr         error
}

func (m *reportUsecaseMock) ProcessInputPaths(paths []string) (domain.ProcessedInput, error) {
	if m.processErr != nil {
		return domain.ProcessedInput{}, m.processErr
	}
	return m.processResult, nil
}

func (m *reportUsecaseMock) GenerateReport(request domain.ReportRequest) (domain.ReportResult, error) {
	m.lastGenerateRequest = request
	if m.generateErr != nil {
		return domain.ReportResult{}, m.generateErr
	}
	return m.generateResult, nil
}

func TestProcessUploadsAggregatesState(t *testing.T) {
	mock := &reportUsecaseMock{
		processResult: domain.ProcessedInput{
			ImagePaths:   []string{"a.png", "b.jpg"},
			IgnoredPaths: []string{"skip.txt"},
		},
	}
	app := newAppWithUsecase(mock)

	got, err := app.ProcessUploads([]string{"a.png", "b.jpg", "skip.txt"})
	if err != nil {
		t.Fatalf("process uploads error: %v", err)
	}

	if got.TotalUploads != 2 {
		t.Fatalf("expected total uploads 2, got %d", got.TotalUploads)
	}
	if len(got.Added) != 2 || len(got.Ignored) != 1 {
		t.Fatalf("unexpected result: %+v", got)
	}
	if len(app.imagePaths) != 2 {
		t.Fatalf("expected app image paths updated, got %d", len(app.imagePaths))
	}
}

func TestRunGenerationUsesAppStateAndReturnsResult(t *testing.T) {
	mock := &reportUsecaseMock{
		generateResult: domain.ReportResult{
			OutputPath:     "relatorio_gerado.docx",
			TotalImages:    3,
			LandscapeCount: 2,
			PortraitCount:  1,
		},
	}
	app := newAppWithUsecase(mock)
	app.templatePath = "template.docx"
	app.imagePaths = []string{"1.png", "2.jpg", "3.jpeg"}

	got, err := app.RunGeneration()
	if err != nil {
		t.Fatalf("run generation error: %v", err)
	}

	if mock.lastGenerateRequest.TemplatePath != "template.docx" {
		t.Fatalf("unexpected template path: %s", mock.lastGenerateRequest.TemplatePath)
	}
	if len(mock.lastGenerateRequest.ImagePaths) != 3 {
		t.Fatalf("expected 3 image paths, got %d", len(mock.lastGenerateRequest.ImagePaths))
	}
	if mock.lastGenerateRequest.OutputPath != "relatorio_gerado.docx" {
		t.Fatalf("unexpected output path: %s", mock.lastGenerateRequest.OutputPath)
	}

	if got.TotalImages != 3 || got.LandscapeCount != 2 || got.PortraitCount != 1 {
		t.Fatalf("unexpected generation result: %+v", got)
	}
}

func TestProcessUploadsAndRunGenerationPropagateErrors(t *testing.T) {
	processErrMock := &reportUsecaseMock{processErr: errors.New("process fail")}
	app := newAppWithUsecase(processErrMock)
	if _, err := app.ProcessUploads([]string{"a.zip"}); err == nil {
		t.Fatal("expected process error")
	}

	generateErrMock := &reportUsecaseMock{generateErr: errors.New("generate fail")}
	app2 := newAppWithUsecase(generateErrMock)
	if _, err := app2.RunGeneration(); err == nil {
		t.Fatal("expected generation error")
	}
}
