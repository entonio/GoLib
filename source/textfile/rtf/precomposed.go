package rtf

import (
	"fmt"

	rtflib "github.com/J45k4/rtf"
)

func Read(rtf string) Rtf {
	return PrecomposedRtf{rtf: rtf}
}

type PrecomposedRtf struct {
	rtf   string
	plain *string
}

func (self PrecomposedRtf) PlainText() (string, error) {
	if self.plain == nil {
		plain, err := rtfToPlain(self.rtf)
		if err != nil {
			return "", err
		}

		self.plain = &plain
	}
	return *self.plain, nil
}

func rtfToPlain(rtf string) (plain string, err error) {
	defer func() {
		if message := recover(); message != nil {
			err = fmt.Errorf("%s", message)
		}
	}()
	plain = rtflib.StripRichTextFormat(rtf)
	return
}
