package charsets

import (
	"golang.org/x/text/encoding/charmap"
)

func Decode(encoded []byte, charset *charmap.Charmap) (decoded string, err error) {
	charsetDecoder := charset.NewDecoder()
	decodedBytes, err := charsetDecoder.Bytes(encoded)
	if err != nil {
		return
	}
	decoded = string(decodedBytes)
	return
}
