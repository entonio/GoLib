package textfile

import (
	"fmt"
	"strings"

	"golib/textfile/rtf"
	"golib/textfile/wordml"
)

type TextFile interface {
	PlainText() (string, error)
}

type Plain string

func (self Plain) PlainText() (string, error) {
	return string(self), nil
}

func Read(s string) TextFile {
	if len(s) == 0 {
		return Plain("")
	}
	if strings.HasPrefix(s, "<?xml") {
		return wordml.Read(s)
	}
	if strings.HasPrefix(s, "{\\rtf") {
		return rtf.Read(s)
	}
	panic(fmt.Sprintf("Could not read text document: %.20s", s))
}
