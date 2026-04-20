package main

import (
	"context"
	"path/filepath"
	"strings"

	"ImageReportGen/core/domain"
	"ImageReportGen/core/service"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type reportUsecase interface {
	ProcessInputPaths(paths []string) (domain.ProcessedInput, error)
	GenerateReport(request domain.ReportRequest) (domain.ReportResult, error)
}

// App struct
type App struct {
	ctx          context.Context
	reportSvc    reportUsecase
	templatePath string
	imagePaths   []string
}

// NewApp creates a new App application struct
func NewApp(svc *service.ReportService) *App {
	return &App{
		reportSvc:  svc,
		imagePaths: []string{},
	}
}

func newAppWithUsecase(svc reportUsecase) *App {
	return &App{
		reportSvc:  svc,
		imagePaths: []string{},
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

type ProcessUploadsResult struct {
	Added        []string `json:"added"`
	Ignored      []string `json:"ignored"`
	TotalUploads int      `json:"totalUploads"`
}

type GenerationResult struct {
	OutputPath     string `json:"outputPath"`
	TotalImages    int    `json:"totalImages"`
	LandscapeCount int    `json:"landscapeCount"`
	PortraitCount  int    `json:"portraitCount"`
}

func (a *App) SelectTemplate() (string, error) {
	path, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Selecionar Modelo .docx",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Documentos Word (*.docx)",
				Pattern:     "*.docx",
			},
		},
	})
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(path) == "" {
		return "", nil
	}
	a.templatePath = path
	return filepath.Base(path), nil
}

func (a *App) ProcessUploads(paths []string) (ProcessUploadsResult, error) {
	processed, err := a.reportSvc.ProcessInputPaths(paths)
	if err != nil {
		return ProcessUploadsResult{}, err
	}

	a.imagePaths = append(a.imagePaths, processed.ImagePaths...)

	return ProcessUploadsResult{
		Added:        processed.ImagePaths,
		Ignored:      processed.IgnoredPaths,
		TotalUploads: len(a.imagePaths),
	}, nil
}

func (a *App) RunGeneration() (GenerationResult, error) {
	result, err := a.reportSvc.GenerateReport(domain.ReportRequest{
		TemplatePath: a.templatePath,
		ImagePaths:   a.imagePaths,
		OutputPath:   "relatorio_gerado.docx",
	})
	if err != nil {
		return GenerationResult{}, err
	}

	return GenerationResult{
		OutputPath:     result.OutputPath,
		TotalImages:    result.TotalImages,
		LandscapeCount: result.LandscapeCount,
		PortraitCount:  result.PortraitCount,
	}, nil
}
