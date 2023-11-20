// Adapted from https://github.com/anansii/docx-templater
package docx

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	"golib/log"

	"golang.org/x/text/unicode/norm"
)

type docxAccessoryFile struct {
	Name     string
	Contents []byte
}

type Docx struct {
	name           string
	accessoryFiles []docxAccessoryFile

	XML     string
	Header2 string
}

/*
func DocxString(original string) string {
	return strings.ReplaceAll(original, "\n", "</w:t></w:r></w:p><w:p><w:r><w:t>")
}
*/

func LoadDocxZ(source string) (*Docx, error) {
	log.Debug("Loading %s", source)
	var doc = &Docx{}
	reader, err := zip.OpenReader(source)
	if err != nil {
		return doc, err
	}

	doc.name = filepath.Base(source)
	doc.XML = ""
	found := false

	for _, file := range reader.File {
		var f docxAccessoryFile
		f.Name = file.Name

		fileReader, err := file.Open()
		if err != nil {
			return doc, err
		}
		defer fileReader.Close()

		bytes, err := io.ReadAll(fileReader)
		if err != nil {
			return doc, err
		}
		f.Contents = bytes

		if f.Name == "word/document.xml" {
			doc.XML = string(f.Contents)
			found = true
		} else if f.Name == "word/header2.xml" {
			doc.Header2 = string(f.Contents)
		} else {
			doc.accessoryFiles = append(doc.accessoryFiles, f)
		}
	}
	defer reader.Close()

	if !found {
		log.Debug("Word XML not found!")
	}

	return doc, nil
}

func (doc *Docx) Save(target string) error {
	docxFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer docxFile.Close()

	writer := zip.NewWriter(docxFile)
	defer writer.Close()

	for _, f := range doc.accessoryFiles {
		zippedFile, err := writer.Create(f.Name)
		if err != nil {
			return err
		}
		zippedFile.Write(f.Contents)
	}

	header2, err := writer.Create("word/header2.xml")
	if err != nil {
		return err
	}
	header2.Write([]byte(norm.NFC.String(doc.Header2)))

	xml, err := writer.Create("word/document.xml")
	if err != nil {
		return err
	}
	xml.Write([]byte(norm.NFC.String(doc.XML)))

	return nil
}
