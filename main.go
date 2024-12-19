package main

import (
	"ImageReportGen/utils"
	"fmt"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Gerador de Relatório de Imagens")

	// Variáveis
	var imagePaths []string
	var templatePath string
	var zipOutputDir = "./temp_images"

	// Labels para exibir informações sobre arquivos selecionados
	templateLabel := widget.NewLabel("Nenhum modelo selecionado")
	imageLabel := widget.NewLabel("Nenhuma imagem selecionada")

	// Lista para exibir as imagens
	imageList := widget.NewList(
		func() int {
			return len(imagePaths)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Imagem")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(imagePaths[i])
		},
	)

	// Botão para selecionar o modelo
	selectTemplateBtn := widget.NewButton("Selecionar Modelo .docx", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				templatePath = reader.URI().Path()
				templateLabel.SetText(fmt.Sprintf("Modelo: %s", filepath.Base(templatePath)))
			}
		}, myWindow)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".docx"}))
		fd.Show()
	})

	// Configura a funcionalidade de drag and drop na janela principal
	myWindow.SetOnDropped(func(position fyne.Position, uris []fyne.URI) {
		for _, uri := range uris {
			filePath := uri.Path()
			if strings.HasSuffix(strings.ToLower(filePath), ".zip") {
				extracted, err := utils.ExtractZip(filePath, zipOutputDir)
				if err != nil {
					dialog.ShowError(err, myWindow)
					continue
				}
				imagePaths = append(imagePaths, extracted...)
			} else if strings.HasSuffix(strings.ToLower(filePath), ".jpg") ||
				strings.HasSuffix(strings.ToLower(filePath), ".jpeg") ||
				strings.HasSuffix(strings.ToLower(filePath), ".png") {
				imagePaths = append(imagePaths, filePath)
			} else {
				dialog.ShowInformation("Arquivo Ignorado", fmt.Sprintf("O arquivo %s não é suportado", filePath), myWindow)
			}
		}

		if len(imagePaths) > 0 {
			imageLabel.SetText(fmt.Sprintf("Imagens Selecionadas: %d", len(imagePaths)))
		} else {
			imageLabel.SetText("Nenhuma imagem selecionada")
		}
		imageList.Refresh()
	})

	// Botão para selecionar imagens
	selectImagesBtn := widget.NewButton("Selecionar Imagens", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				imagePaths = append(imagePaths, reader.URI().Path())
				imageList.Refresh()
				imageLabel.SetText(fmt.Sprintf("Imagens Selecionadas: %d", len(imagePaths)))
			}
		}, myWindow)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".jpg", ".jpeg"}))
		fd.Show()
	})


	// Botão para gerar o relatório
	generateReportBtn := widget.NewButton("Gerar Relatório", func() {
		defer utils.CleanUpTempFiles(zipOutputDir) // Limpar arquivos temporários
		if len(imagePaths) == 0 {
			dialog.ShowInformation("Erro", "Por favor, adicione pelo menos uma imagem", myWindow)
			return
		}

		err := utils.GenerateReport(imagePaths,templatePath)
		if err != nil {
			dialog.ShowError(err, myWindow)
		} else {
			dialog.ShowInformation("Sucesso", "Relatório gerado com sucesso!", myWindow)
		}
	})

	// Layout com espaçamento
	content := container.NewVBox(
		templateLabel,
		selectTemplateBtn,
		widget.NewSeparator(), // Separador para organização visual
		widget.NewLabel("Arraste arquivos de imagem ou um .zip para a janela:"),
		selectImagesBtn,
		imageLabel,
		imageList,
		generateReportBtn,
	)

	// Ajuste da janela
	myWindow.SetContent(container.NewPadded(content)) // Adiciona padding ao layout
	myWindow.Resize(fyne.NewSize(600, 500))
	myWindow.ShowAndRun()
}
