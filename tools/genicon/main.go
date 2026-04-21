// genicon builds a valid Windows .ico from a PNG source (multiple embedded PNG sizes).
package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"image"
	"image/png"
	"os"
)

func main() {
	inPath := flag.String("in", "build/appicon.png", "input PNG")
	outPath := flag.String("out", "build/windows/icon.ico", "output ICO")
	flag.Parse()

	srcData, err := os.ReadFile(*inPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	srcImg, err := png.Decode(bytes.NewReader(srcData))
	if err != nil {
		fmt.Fprintln(os.Stderr, "decode png:", err)
		os.Exit(1)
	}

	sizes := []int{256, 128, 64, 48, 32, 16}
	var pngChunks [][]byte
	for _, sz := range sizes {
		resized := resizeNearest(srcImg, sz, sz)
		var buf bytes.Buffer
		if err := png.Encode(&buf, resized); err != nil {
			fmt.Fprintln(os.Stderr, "encode png:", err)
			os.Exit(1)
		}
		pngChunks = append(pngChunks, buf.Bytes())
	}

	if err := writeICO(*outPath, pngChunks); err != nil {
		fmt.Fprintln(os.Stderr, "write ico:", err)
		os.Exit(1)
	}
	fmt.Println("wrote", *outPath)
}

func resizeNearest(src image.Image, w, h int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	sb := src.Bounds()
	sw := sb.Dx()
	sh := sb.Dy()
	if sw == 0 || sh == 0 {
		return dst
	}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			sx := sb.Min.X + (x*sw)/w
			sy := sb.Min.Y + (y*sh)/h
			dst.Set(x, y, src.At(sx, sy))
		}
	}
	return dst
}

// writeICO writes an ICO containing PNG-encoded images (Windows Vista+).
func writeICO(path string, pngImages [][]byte) error {
	n := len(pngImages)
	if n == 0 {
		return fmt.Errorf("no images")
	}
	headerSize := 6 + 16*n
	offset := uint32(headerSize)
	var dir [6]byte
	binary.LittleEndian.PutUint16(dir[0:2], 0) // reserved
	binary.LittleEndian.PutUint16(dir[2:4], 1) // type: icon
	binary.LittleEndian.PutUint16(dir[4:6], uint16(n))

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(dir[:]); err != nil {
		return err
	}

	for _, pngData := range pngImages {
		w, h, err := pngDimensions(pngData)
		if err != nil {
			return err
		}
		var entry [16]byte
		bw, bh := byte(w), byte(h)
		if w >= 256 {
			bw = 0
		}
		if h >= 256 {
			bh = 0
		}
		entry[0] = bw
		entry[1] = bh
		entry[2] = 0 // color count
		entry[3] = 0 // reserved
		binary.LittleEndian.PutUint16(entry[4:6], 1) // planes
		binary.LittleEndian.PutUint16(entry[6:8], 32)
		binary.LittleEndian.PutUint32(entry[8:12], uint32(len(pngData)))
		binary.LittleEndian.PutUint32(entry[12:16], offset)
		if _, err := f.Write(entry[:]); err != nil {
			return err
		}
		offset += uint32(len(pngData))
	}

	for _, pngData := range pngImages {
		if _, err := f.Write(pngData); err != nil {
			return err
		}
	}
	return nil
}

func pngDimensions(pngData []byte) (int, int, error) {
	cfg, err := png.DecodeConfig(bytes.NewReader(pngData))
	if err != nil {
		return 0, 0, err
	}
	return cfg.Width, cfg.Height, nil
}
