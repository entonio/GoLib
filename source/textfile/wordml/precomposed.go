package wordml

import (
	"regexp"
	"strings"
)

func Read(xml string) WordML {
	return PrecomposedWordML{xml: xml}
}

type PrecomposedWordML struct {
	xml   string
	plain *string
}

func (self PrecomposedWordML) PlainText() (string, error) {
	if self.plain == nil {
		plain := xmlToPlain(self.xml)

		self.plain = &plain
	}
	return *self.plain, nil
}

func xmlToPlain(xml string) string {
	var d _document
	d.listP(xml)
	var document strings.Builder
	for _, paragraph := range d.paragraphs {
		for _, run := range paragraph.runs {
			document.WriteString(run)
		}
		document.WriteString("\n")
	}
	return document.String()
}

type _document struct {
	paragraphs []*_paragraph
}

type _paragraph struct {
	runs []string
}

var reParagraph = regexp.MustCompile(`(?U)<w:p[^>]*(.*)</w:p>`)
var reRun = regexp.MustCompile(`(?U)(<w:r>|<w:r .*>)(.*)(</w:r>)`)
var reContent = regexp.MustCompile(`(?U)(<w:t>|<w:t .*>)(.*)(</w:t>)`)

// listP for w:p tag value
func (d *_document) listP(data string) {
	var result []string
	for _, match := range reParagraph.FindAllStringSubmatch(data, -1) {
		result = append(result, match[1])
	}
	for _, item := range result {
		if hasP(item) {
			d.listP(item)
			continue
		}
		d.getT(item)
	}
}

// hasP identify the paragraph
func hasP(data string) bool {
	re := regexp.MustCompile(`(?U)<w:p[^>]*>(.*)</w:p>`)
	result := re.MatchString(data)
	return result
}

// get w:t value
func (d *_document) getT(data string) {
	var subStr string
	content := []string{}
	wrMatch := reRun.FindAllStringSubmatchIndex(data, -1)
	// loop r
	for _, rMatch := range wrMatch {
		rData := data[rMatch[4]:rMatch[5]]
		wtMatch := reContent.FindAllStringSubmatchIndex(rData, -1)
		for _, match := range wtMatch {
			subStr = rData[match[4]:match[5]]
			content = append(content, subStr)
		}
	}

	w := new(_paragraph)
	w.runs = content
	d.paragraphs = append(d.paragraphs, w)
}
