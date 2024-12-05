package utils

import (
	"os"

	"github.com/gomutex/godocx"
	"github.com/gomutex/godocx/common/units"
	"github.com/gomutex/godocx/docx"
)

func GenerateReport(imagePaths []string, templatePath string) error {
	doc, err := godocx.OpenDocument(templatePath)
	if err != nil {
		return err
	}
	defer doc.Close()

	var landscapeImages []string
	var portraitImages []string

	for _, imgPath := range imagePaths {
		orientation, err := getImageOrientation(imgPath)
		if err != nil {
			return err
		}

		if orientation == "Paisagem" {
			landscapeImages = append(landscapeImages, imgPath)
		} else {
			portraitImages = append(portraitImages, imgPath)
		}
	}

	// Adicionar imagens em paisagem
	if len(landscapeImages) > 0 {
		err = addImagesToDocument(doc, landscapeImages, 3, "Paisagem")
		if err != nil {
			return err
		}
		doc.AddPageBreak()
	}

	// Adicionar imagens em retrato
	if len(portraitImages) > 0 {
		err = addImagesToDocument(doc, portraitImages, 4, "Retrato")
		if err != nil {
			return err
		}
		doc.AddPageBreak()
	}

	// Salvar o documento
	outputPath := "relatorio_gerado.docx"
	err = doc.SaveTo(outputPath)
	if err != nil {
		return err
	}

	return nil
}

func addImagesToDocument(doc *docx.RootDoc, images []string, imagesPerPage int, orientation string) error {
	imageCount := 0
	totalImages := len(images)

	for imageCount < totalImages {
		// Adicionar uma quebra de página, se não for a primeira página
		if imageCount > 0 {
			doc.AddPageBreak()
		}

		// Criar uma tabela para organizar as imagens
		var rows, cols int
		var width, height float64
		if orientation == "Paisagem" {
			rows = 1
			cols = 3 // Até 3 imagens por linha para paisagem
			width = 5
			height = 3
		} else {
			rows = 2
			cols = 2 // Até 4 imagens no total, 2x2 para retrato
			width = 3
			height = 5
		}

		table := doc.AddTable()

		for r := 0; r < rows && imageCount < totalImages; r++ {
			row := table.AddRow()
			for c := 0; c < cols && imageCount < totalImages; c++ {
				cell := row.AddCell()
				err := insertImageIntoCell(cell, images[imageCount], width, height)
				if err != nil {
					return err
				}
				imageCount++
			}
		}
	}

	return nil
}

func insertImageIntoCell(cell *docx.Cell, imagePath string, width, height float64) error {

	// Criar um novo parágrafo na célula
	_, err := cell.AddEmptyPara().AddPicture(imagePath, units.Inch(width), units.Inch(height))
	if err != nil {
		return err
	}

	return nil
}

func CleanUpTempFiles(tempDir string) {
	os.RemoveAll(tempDir)
}
