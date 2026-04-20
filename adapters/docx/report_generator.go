package docx

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const (
	documentXMLPath     = "word/document.xml"
	documentRelsXMLPath = "word/_rels/document.xml.rels"
	contentTypesXMLPath = "[Content_Types].xml"
	imageRelType        = "http://schemas.openxmlformats.org/officeDocument/2006/relationships/image"
	imageTargetPrefix   = "media/"
	emuPerInch          = 914400
)

type ReportGenerator struct{}

type pendingImage struct {
	sourcePath string
	zipPath    string
	relID      string
}

type relationships struct {
	XMLName       xml.Name       `xml:"Relationships"`
	XMLNS         string         `xml:"xmlns,attr,omitempty"`
	Relationships []relationship `xml:"Relationship"`
}

type relationship struct {
	ID     string `xml:"Id,attr"`
	Type   string `xml:"Type,attr"`
	Target string `xml:"Target,attr"`
}

type contentTypes struct {
	XMLName   xml.Name    `xml:"Types"`
	XMLNS     string      `xml:"xmlns,attr,omitempty"`
	Defaults  []ctDefault `xml:"Default"`
	Overrides []ctAnyPart `xml:"Override"`
}

type ctDefault struct {
	Extension   string `xml:"Extension,attr"`
	ContentType string `xml:"ContentType,attr"`
}

type ctAnyPart struct {
	PartName    string `xml:"PartName,attr"`
	ContentType string `xml:"ContentType,attr"`
}

func NewReportGenerator() *ReportGenerator {
	return &ReportGenerator{}
}

func (g *ReportGenerator) GenerateReport(templatePath string, landscapeImages, portraitImages []string, outputPath string) error {
	template, err := zip.OpenReader(templatePath)
	if err != nil {
		return err
	}
	defer template.Close()

	entries := map[string][]byte{}
	existingMedia := map[string]struct{}{}

	for _, file := range template.File {
		data, err := readZipFile(file)
		if err != nil {
			return err
		}
		entries[file.Name] = data
		if strings.HasPrefix(file.Name, "word/media/") {
			existingMedia[strings.TrimPrefix(file.Name, "word/media/")] = struct{}{}
		}
	}

	relsData, ok := entries[documentRelsXMLPath]
	if !ok {
		return fmt.Errorf("missing %s", documentRelsXMLPath)
	}
	rels, err := parseRelationships(relsData)
	if err != nil {
		return err
	}

	allImages := append(append([]string{}, landscapeImages...), portraitImages...)
	newImages := make([]pendingImage, 0, len(allImages))

	for _, imagePath := range allImages {
		mediaName, err := nextMediaName(existingMedia, imagePath)
		if err != nil {
			return err
		}
		existingMedia[mediaName] = struct{}{}

		relID := nextRelID(rels)
		rels.Relationships = append(rels.Relationships, relationship{
			ID:     relID,
			Type:   imageRelType,
			Target: imageTargetPrefix + mediaName,
		})

		newImages = append(newImages, pendingImage{
			sourcePath: imagePath,
			zipPath:    "word/media/" + mediaName,
			relID:      relID,
		})
	}

	imageRelIDs := make([]string, 0, len(newImages))
	for _, img := range newImages {
		imageRelIDs = append(imageRelIDs, img.relID)
	}
	landscapeRelIDs := imageRelIDs[:len(landscapeImages)]
	portraitRelIDs := imageRelIDs[len(landscapeImages):]

	documentData, ok := entries[documentXMLPath]
	if !ok {
		return fmt.Errorf("missing %s", documentXMLPath)
	}
	newDocumentXML, err := injectDocumentBody(documentData, buildGeneratedBodyXML(landscapeRelIDs, portraitRelIDs))
	if err != nil {
		return err
	}
	entries[documentXMLPath] = newDocumentXML

	newRelsXML, err := marshalRelationships(rels)
	if err != nil {
		return err
	}
	entries[documentRelsXMLPath] = newRelsXML

	contentTypesData, ok := entries[contentTypesXMLPath]
	if !ok {
		return fmt.Errorf("missing %s", contentTypesXMLPath)
	}
	newContentTypes, err := ensureImageContentTypes(contentTypesData, newImages)
	if err != nil {
		return err
	}
	entries[contentTypesXMLPath] = newContentTypes

	for _, img := range newImages {
		data, err := os.ReadFile(img.sourcePath)
		if err != nil {
			return err
		}
		entries[img.zipPath] = data
	}

	return writeDocx(entries, outputPath)
}

func readZipFile(file *zip.File) ([]byte, error) {
	rc, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return io.ReadAll(rc)
}

func parseRelationships(data []byte) (relationships, error) {
	var rels relationships
	if err := xml.Unmarshal(data, &rels); err != nil {
		return relationships{}, err
	}
	if rels.XMLNS == "" {
		rels.XMLNS = "http://schemas.openxmlformats.org/package/2006/relationships"
	}
	return rels, nil
}

func marshalRelationships(rels relationships) ([]byte, error) {
	out, err := xml.Marshal(rels)
	if err != nil {
		return nil, err
	}
	return append([]byte(xml.Header), out...), nil
}

func nextRelID(rels relationships) string {
	maxID := 0
	for _, rel := range rels.Relationships {
		if strings.HasPrefix(rel.ID, "rId") {
			n, err := strconv.Atoi(strings.TrimPrefix(rel.ID, "rId"))
			if err == nil && n > maxID {
				maxID = n
			}
		}
	}
	return fmt.Sprintf("rId%d", maxID+1)
}

func nextMediaName(existing map[string]struct{}, imagePath string) (string, error) {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(imagePath), "."))
	if ext == "" {
		return "", errors.New("image file without extension")
	}
	for i := 1; ; i++ {
		candidate := fmt.Sprintf("image%d.%s", i, ext)
		if _, exists := existing[candidate]; !exists {
			return candidate, nil
		}
	}
}

func injectDocumentBody(documentXML []byte, generatedXML string) ([]byte, error) {
	raw := string(documentXML)
	bodyOpenStart := strings.Index(raw, "<w:body")
	if bodyOpenStart == -1 {
		return nil, errors.New("document.xml missing w:body start")
	}
	bodyOpenEnd := strings.Index(raw[bodyOpenStart:], ">")
	if bodyOpenEnd == -1 {
		return nil, errors.New("document.xml malformed w:body start tag")
	}
	bodyOpenEnd += bodyOpenStart

	bodyCloseStart := strings.Index(raw, "</w:body>")
	if bodyCloseStart == -1 || bodyCloseStart < bodyOpenEnd {
		return nil, errors.New("document.xml missing w:body end")
	}

	bodyInner := raw[bodyOpenEnd+1 : bodyCloseStart]
	sectStart := strings.LastIndex(bodyInner, "<w:sectPr")
	if sectStart == -1 {
		return nil, errors.New("document.xml missing w:sectPr")
	}
	sectEnd := strings.Index(bodyInner[sectStart:], "</w:sectPr>")
	if sectEnd == -1 {
		return nil, errors.New("document.xml malformed w:sectPr")
	}
	sectEnd += sectStart + len("</w:sectPr>")

	beforeSect := strings.TrimRight(bodyInner[:sectStart], " \r\n\t")
	sectBlock := bodyInner[sectStart:sectEnd]
	afterSect := strings.TrimSpace(bodyInner[sectEnd:])
	if afterSect != "" {
		return nil, errors.New("unexpected body content after w:sectPr")
	}

	newBodyInner := beforeSect
	if generatedXML != "" {
		newBodyInner += generatedXML
	}
	newBodyInner += sectBlock

	result := raw[:bodyOpenEnd+1] + newBodyInner + raw[bodyCloseStart:]
	return []byte(result), nil
}

func buildGeneratedBodyXML(landscapeRelIDs, portraitRelIDs []string) string {
	var builder strings.Builder
	if len(landscapeRelIDs) > 0 {
		writeOrientationTables(&builder, landscapeRelIDs, "Paisagem")
		builder.WriteString(pageBreakXML())
	}
	if len(portraitRelIDs) > 0 {
		writeOrientationTables(&builder, portraitRelIDs, "Retrato")
		builder.WriteString(pageBreakXML())
	}
	return builder.String()
}

func writeOrientationTables(builder *strings.Builder, relIDs []string, orientation string) {
	rows, cols, widthIn, heightIn := layoutByOrientation(orientation)
	idx := 0
	for idx < len(relIDs) {
		if idx > 0 {
			builder.WriteString(pageBreakXML())
		}
		builder.WriteString("<w:tbl><w:tblPr/><w:tblGrid/>")
		for r := 0; r < rows && idx < len(relIDs); r++ {
			builder.WriteString("<w:tr>")
			for c := 0; c < cols; c++ {
				builder.WriteString("<w:tc>")
				if idx < len(relIDs) {
					builder.WriteString(imageParagraphXML(relIDs[idx], widthIn, heightIn))
					idx++
				} else {
					builder.WriteString("<w:p/>")
				}
				builder.WriteString("</w:tc>")
			}
			builder.WriteString("</w:tr>")
		}
		builder.WriteString("</w:tbl>")
	}
}

func imageParagraphXML(relID string, widthInch, heightInch float64) string {
	cx := int(widthInch * emuPerInch)
	cy := int(heightInch * emuPerInch)
	return fmt.Sprintf(
		`<w:p><w:r><w:drawing><wp:inline distT="0" distB="0" distL="0" distR="0"><wp:extent cx="%d" cy="%d"/><wp:docPr id="1" name="Picture"/><a:graphic><a:graphicData uri="http://schemas.openxmlformats.org/drawingml/2006/picture"><pic:pic><pic:nvPicPr><pic:cNvPr id="0" name="image"/><pic:cNvPicPr/></pic:nvPicPr><pic:blipFill><a:blip r:embed="%s"/><a:stretch><a:fillRect/></a:stretch></pic:blipFill><pic:spPr><a:xfrm><a:off x="0" y="0"/><a:ext cx="%d" cy="%d"/></a:xfrm><a:prstGeom prst="rect"><a:avLst/></a:prstGeom></pic:spPr></pic:pic></a:graphicData></a:graphic></wp:inline></w:drawing></w:r></w:p>`,
		cx, cy, relID, cx, cy,
	)
}

func pageBreakXML() string {
	return `<w:p><w:r><w:br w:type="page"/></w:r></w:p>`
}

func layoutByOrientation(orientation string) (int, int, float64, float64) {
	if orientation == "Paisagem" {
		return 3, 1, 4, 2.5
	}
	return 2, 2, 3, 4
}

func ensureImageContentTypes(data []byte, images []pendingImage) ([]byte, error) {
	var ct contentTypes
	if err := xml.Unmarshal(data, &ct); err != nil {
		return nil, err
	}
	if ct.XMLNS == "" {
		ct.XMLNS = "http://schemas.openxmlformats.org/package/2006/content-types"
	}

	existing := map[string]struct{}{}
	for _, def := range ct.Defaults {
		existing[strings.ToLower(def.Extension)] = struct{}{}
	}

	for _, img := range images {
		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(img.sourcePath), "."))
		if ext == "" {
			continue
		}
		if _, ok := existing[ext]; ok {
			continue
		}
		contentType := "image/" + ext
		if ext == "jpg" {
			contentType = "image/jpeg"
		}
		ct.Defaults = append(ct.Defaults, ctDefault{
			Extension:   ext,
			ContentType: contentType,
		})
		existing[ext] = struct{}{}
	}

	out, err := xml.Marshal(ct)
	if err != nil {
		return nil, err
	}
	return append([]byte(xml.Header), out...), nil
}

func writeDocx(entries map[string][]byte, outputPath string) error {
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	zipWriter := zip.NewWriter(outFile)
	defer zipWriter.Close()

	names := make([]string, 0, len(entries))
	for name := range entries {
		names = append(names, name)
	}
	sort.Strings(names)

	for _, name := range names {
		w, err := zipWriter.Create(name)
		if err != nil {
			return err
		}
		if _, err := io.Copy(w, bytes.NewReader(entries[name])); err != nil {
			return err
		}
	}

	return zipWriter.Close()
}
