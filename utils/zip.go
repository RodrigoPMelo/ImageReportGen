package utils

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ExtractZip(zipPath string, outputDir string) ([]string, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var extractedFiles []string
	for _, file := range r.File {
		// Validar extensão do arquivo
		if !strings.HasSuffix(strings.ToLower(file.Name), ".jpg") &&
			!strings.HasSuffix(strings.ToLower(file.Name), ".jpeg") &&
			!strings.HasSuffix(strings.ToLower(file.Name), ".png") {
			continue
		}

		// Criar o destino
		destPath := filepath.Join(outputDir, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(destPath, os.ModePerm)
			continue
		}

		// Criar o arquivo
		os.MkdirAll(filepath.Dir(destPath), os.ModePerm)
		destFile, err := os.Create(destPath)
		if err != nil {
			return nil, err
		}
		defer destFile.Close()

		srcFile, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer srcFile.Close()

		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			return nil, err
		}

		extractedFiles = append(extractedFiles, destPath)
	}

	return extractedFiles, nil
}
