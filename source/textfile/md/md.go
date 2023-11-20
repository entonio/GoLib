//go:build descontinuado

package md

import (
	"fmt"
	"strings"
)

type Paragraph struct {
	Alignment Alignment
	Runs      []Run
}

type Alignment string

const (
	AlignCenter  Alignment = "c"
	AlignLeft    Alignment = "l"
	AlignRight   Alignment = "r"
	AlignJustify Alignment = "j"
)

type Run struct {
	Text   string
	IsBold bool
}

func Parse(s string) (paragraphs []Paragraph) {
	for _, p := range strings.Split(s, "\n") {
		paragraphs = append(paragraphs, NewParagraph(p))
	}
	return
}

func NewParagraph(s string) (paragraph Paragraph) {
	content := strings.TrimPrefix(s, "/c:")
	if len(content) != len(s) {
		paragraph.Alignment = AlignCenter
	} else {
		paragraph.Alignment = AlignJustify
	}

	run := Run{}
	for i, c := range content {
		if c != '*' {
			run.Text += string(c)
			continue
		}
		if len(run.Text) == 0 {
			run.IsBold = !run.IsBold
			fmt.Printf("Toggling bold %t for %s\n", run.IsBold, content[i:])
			continue
		}
		previous := run
		paragraph.Runs = append(paragraph.Runs, previous)
		run = Run{}
		run.IsBold = !previous.IsBold
		fmt.Printf("Closing %v, starting bold %t for %s\n", previous, run.IsBold, content[i:])
	}

	if len(run.Text) > 0 {
		fmt.Printf("Closing %v\n", run)
		paragraph.Runs = append(paragraph.Runs, run)
	} else {
		fmt.Printf("No pending close")
	}

	return
}
