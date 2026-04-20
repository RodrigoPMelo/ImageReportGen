package docx

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateReportPreservesSectPrAndHeaderFooterReferences(t *testing.T) {
	tempDir := t.TempDir()
	templatePath := filepath.Join(tempDir, "template.docx")
	outputPath := filepath.Join(tempDir, "output.docx")

	imgA := filepath.Join(tempDir, "a.png")
	imgB := filepath.Join(tempDir, "b.png")
	imgC := filepath.Join(tempDir, "c.png")
	for _, p := range []string{imgA, imgB, imgC} {
		if err := os.WriteFile(p, samplePNG(), 0o644); err != nil {
			t.Fatalf("write png %s: %v", p, err)
		}
	}

	if err := writeTemplateDocx(templatePath); err != nil {
		t.Fatalf("write template docx: %v", err)
	}

	gen := NewReportGenerator()
	err := gen.GenerateReport(templatePath, []string{imgA, imgB}, []string{imgC}, outputPath)
	if err != nil {
		t.Fatalf("generate report: %v", err)
	}

	files, err := readZipEntries(outputPath)
	if err != nil {
		t.Fatalf("read output zip: %v", err)
	}

	if _, ok := files["word/header1.xml"]; !ok {
		t.Fatal("header1.xml missing from output")
	}
	if _, ok := files["word/footer1.xml"]; !ok {
		t.Fatal("footer1.xml missing from output")
	}

	docXML := string(files["word/document.xml"])
	bodyStart := strings.Index(docXML, "<w:body>")
	bodyEnd := strings.Index(docXML, "</w:body>")
	if bodyStart == -1 || bodyEnd == -1 || bodyEnd <= bodyStart {
		t.Fatal("invalid w:body region")
	}
	bodyInner := docXML[bodyStart+len("<w:body>") : bodyEnd]

	sectStart := strings.LastIndex(bodyInner, "<w:sectPr")
	sectEnd := strings.LastIndex(bodyInner, "</w:sectPr>")
	if sectStart == -1 || sectEnd == -1 || sectEnd < sectStart {
		t.Fatal("w:sectPr block not found in body")
	}
	afterSect := strings.TrimSpace(bodyInner[sectEnd+len("</w:sectPr>"):])
	if afterSect != "" {
		t.Fatalf("w:sectPr is not final element in body, trailing content: %q", afterSect)
	}

	sectBlock := bodyInner[sectStart : sectEnd+len("</w:sectPr>")]
	if !strings.Contains(sectBlock, "headerReference") || !strings.Contains(sectBlock, "footerReference") {
		t.Fatal("w:sectPr lost header/footer references")
	}

	relsXML := string(files["word/_rels/document.xml.rels"])
	if !strings.Contains(relsXML, "header1.xml") || !strings.Contains(relsXML, "footer1.xml") {
		t.Fatal("document rels lost header/footer targets")
	}

	mediaCount := 0
	for name := range files {
		if strings.HasPrefix(name, "word/media/") {
			mediaCount++
		}
	}
	if mediaCount != 3 {
		t.Fatalf("expected 3 media files, got %d", mediaCount)
	}

	if err := xml.Unmarshal(files["word/document.xml"], new(interface{})); err != nil {
		t.Fatalf("document.xml is not well formed: %v", err)
	}
	if err := xml.Unmarshal(files["word/_rels/document.xml.rels"], new(interface{})); err != nil {
		t.Fatalf("document.xml.rels is not well formed: %v", err)
	}
}

func writeTemplateDocx(path string) error {
	buf := bytes.NewBuffer(nil)
	zw := zip.NewWriter(buf)

	entries := map[string]string{
		"[Content_Types].xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
  <Override PartName="/word/header1.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.header+xml"/>
  <Override PartName="/word/footer1.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.footer+xml"/>
</Types>`,
		"_rels/.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>
</Relationships>`,
		"word/document.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"
 xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships"
 xmlns:wp="http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing"
 xmlns:a="http://schemas.openxmlformats.org/drawingml/2006/main"
 xmlns:pic="http://schemas.openxmlformats.org/drawingml/2006/picture">
  <w:body>
    <w:p><w:r><w:t>Template Body</w:t></w:r></w:p>
    <w:sectPr>
      <w:headerReference w:type="default" r:id="rIdHeader1"/>
      <w:footerReference w:type="default" r:id="rIdFooter1"/>
      <w:pgSz w:w="11906" w:h="16838"/>
    </w:sectPr>
  </w:body>
</w:document>`,
		"word/_rels/document.xml.rels": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rIdHeader1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/header" Target="header1.xml"/>
  <Relationship Id="rIdFooter1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/footer" Target="footer1.xml"/>
</Relationships>`,
		"word/header1.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:hdr xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:p><w:r><w:t>Header</w:t></w:r></w:p></w:hdr>`,
		"word/footer1.xml": `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:ftr xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main"><w:p><w:r><w:t>Footer</w:t></w:r></w:p></w:ftr>`,
	}

	for name, content := range entries {
		w, err := zw.Create(name)
		if err != nil {
			return err
		}
		if _, err := w.Write([]byte(content)); err != nil {
			return err
		}
	}

	if err := zw.Close(); err != nil {
		return err
	}
	return os.WriteFile(path, buf.Bytes(), 0o644)
}

func readZipEntries(path string) (map[string][]byte, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	out := make(map[string][]byte, len(r.File))
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		data := bytes.NewBuffer(nil)
		if _, err := data.ReadFrom(rc); err != nil {
			rc.Close()
			return nil, err
		}
		if err := rc.Close(); err != nil {
			return nil, err
		}
		out[f.Name] = data.Bytes()
	}
	return out, nil
}

func samplePNG() []byte {
	return []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A,
		0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4,
		0x89, 0x00, 0x00, 0x00, 0x0D, 0x49, 0x44, 0x41,
		0x54, 0x78, 0x9C, 0x63, 0xF8, 0xCF, 0xC0, 0xF0,
		0x1F, 0x00, 0x05, 0x00, 0x01, 0xFF, 0x89, 0x99,
		0x3D, 0x1D, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45,
		0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82,
	}
}
