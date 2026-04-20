package fs

import (
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"ImageReportGen/core/domain"
)

func TestGetImageOrientationLandscapeAndPortrait(t *testing.T) {
	tmp := t.TempDir()
	landscape := filepath.Join(tmp, "landscape.png")
	portrait := filepath.Join(tmp, "portrait.png")

	if err := writePNG(landscape, 200, 100); err != nil {
		t.Fatalf("write landscape: %v", err)
	}
	if err := writePNG(portrait, 100, 200); err != nil {
		t.Fatalf("write portrait: %v", err)
	}

	p := NewImageProcessor()

	landscapeOrientation, err := p.GetImageOrientation(landscape)
	if err != nil {
		t.Fatalf("landscape orientation error: %v", err)
	}
	if landscapeOrientation != domain.OrientationLandscape {
		t.Fatalf("expected %s, got %s", domain.OrientationLandscape, landscapeOrientation)
	}

	portraitOrientation, err := p.GetImageOrientation(portrait)
	if err != nil {
		t.Fatalf("portrait orientation error: %v", err)
	}
	if portraitOrientation != domain.OrientationPortrait {
		t.Fatalf("expected %s, got %s", domain.OrientationPortrait, portraitOrientation)
	}
}

func TestGetImageOrientationReturnsErrorForInvalidOrMissingFile(t *testing.T) {
	tmp := t.TempDir()
	invalid := filepath.Join(tmp, "invalid.png")
	if err := os.WriteFile(invalid, []byte("not-image-content"), 0o644); err != nil {
		t.Fatalf("write invalid: %v", err)
	}

	p := NewImageProcessor()

	if _, err := p.GetImageOrientation(filepath.Join(tmp, "missing.png")); err == nil {
		t.Fatal("expected error for missing file")
	}

	if _, err := p.GetImageOrientation(invalid); err == nil {
		t.Fatal("expected error for invalid image file")
	}
}

func writePNG(path string, width, height int) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{R: 100, G: 100, B: 100, A: 255})
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}
