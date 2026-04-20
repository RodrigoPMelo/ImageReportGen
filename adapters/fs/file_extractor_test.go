package fs

import (
	"archive/zip"
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestExtractZipExtractsOnlySupportedImages(t *testing.T) {
	tmp := t.TempDir()
	zipPath := filepath.Join(tmp, "images.zip")
	outDir := filepath.Join(tmp, "out")

	files := map[string][]byte{
		"one.png":  []byte("png-data"),
		"two.jpg":  []byte("jpg-data"),
		"three.md": []byte("not-image"),
	}
	if err := writeZip(zipPath, files); err != nil {
		t.Fatalf("write zip: %v", err)
	}

	extractor := NewFileExtractor()
	got, err := extractor.ExtractZip(zipPath, outDir)
	if err != nil {
		t.Fatalf("extract zip: %v", err)
	}

	sort.Strings(got)
	want := []string{
		filepath.Join(outDir, "one.png"),
		filepath.Join(outDir, "two.jpg"),
	}
	sort.Strings(want)

	if len(got) != len(want) {
		t.Fatalf("expected %d files, got %d (%v)", len(want), len(got), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected path %s, got %s", want[i], got[i])
		}
	}

	if _, err := os.Stat(filepath.Join(outDir, "three.md")); err == nil {
		t.Fatal("unsupported extension should not be extracted")
	}
}

func TestExtractZipInvalidZipReturnsError(t *testing.T) {
	tmp := t.TempDir()
	invalidZip := filepath.Join(tmp, "invalid.zip")
	if err := os.WriteFile(invalidZip, []byte("not-a-zip"), 0o644); err != nil {
		t.Fatalf("write invalid zip: %v", err)
	}

	extractor := NewFileExtractor()
	if _, err := extractor.ExtractZip(invalidZip, filepath.Join(tmp, "out")); err == nil {
		t.Fatal("expected error for invalid zip")
	}
}

func TestCleanUpTempFilesRemovesDirectory(t *testing.T) {
	tmp := t.TempDir()
	target := filepath.Join(tmp, "cleanup")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(target, "temp.txt"), []byte("x"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	extractor := NewFileExtractor()
	if err := extractor.CleanUpTempFiles(target); err != nil {
		t.Fatalf("cleanup: %v", err)
	}

	if _, err := os.Stat(target); !os.IsNotExist(err) {
		t.Fatalf("expected directory removed, stat err: %v", err)
	}
}

func writeZip(path string, files map[string][]byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	zw := zip.NewWriter(f)
	for name, content := range files {
		w, err := zw.Create(name)
		if err != nil {
			_ = zw.Close()
			return err
		}
		if _, err := w.Write(content); err != nil {
			_ = zw.Close()
			return err
		}
	}
	return zw.Close()
}
