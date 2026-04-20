package main

import (
	"embed"

	docxadapter "ImageReportGen/adapters/docx"
	fsadapter "ImageReportGen/adapters/fs"
	"ImageReportGen/core/service"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	fileExtractor := fsadapter.NewFileExtractor()
	imageProcessor := fsadapter.NewImageProcessor()
	reportGenerator := docxadapter.NewReportGenerator()
	reportService := service.NewReportService(fileExtractor, imageProcessor, reportGenerator, "./temp_images")
	app := NewApp(reportService)

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "ImageReportGen",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		DragAndDrop: &options.DragAndDrop{
			EnableFileDrop:    true,
			DisableWebViewDrop: true,
		},
		Windows: &windows.Options{},
		Mac:     &mac.Options{},
		Linux:   &linux.Options{},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
