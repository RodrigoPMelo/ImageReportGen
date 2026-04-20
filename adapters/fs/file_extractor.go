package fs

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type FileExtractor struct{}

func NewFileExtractor() *FileExtractor {
	return &FileExtractor{}
}

func (f *FileExtractor) ExtractZip(zipPath string, outputDir string) ([]string, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	extractedFiles := []string{}
	for _, file := range r.File {
		if !isSupportedImage(file.Name) {
			continue
		}

		destPath := filepath.Join(outputDir, file.Name)
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(destPath, os.ModePerm); err != nil {
				return nil, err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
			return nil, err
		}

		destFile, err := os.Create(destPath)
		if err != nil {
			return nil, err
		}

		srcFile, err := file.Open()
		if err != nil {
			destFile.Close()
			return nil, err
		}

		_, copyErr := io.Copy(destFile, srcFile)
		closeSrcErr := srcFile.Close()
		closeDestErr := destFile.Close()

		if copyErr != nil {
			return nil, copyErr
		}
		if closeSrcErr != nil {
			return nil, closeSrcErr
		}
		if closeDestErr != nil {
			return nil, closeDestErr
		}

		extractedFiles = append(extractedFiles, destPath)
	}

	return extractedFiles, nil
}

func (f *FileExtractor) CleanUpTempFiles(tempDir string) error {
	return os.RemoveAll(tempDir)
}

func isSupportedImage(path string) bool {
	lower := strings.ToLower(path)
	return strings.HasSuffix(lower, ".jpg") ||
		strings.HasSuffix(lower, ".jpeg") ||
		strings.HasSuffix(lower, ".png")
}
