package dbq

import (
	"bytes"
	"compress/zlib"
	"io"
	"strings"

	"golib/text"
)

func SPlain[T any](query string, label string, values []T) string {
	return strings.ReplaceAll(query, label,
		text.Join("", ",", ",", "", values),
	)
}

func SText[T any](query string, label string, values []T) string {
	return STextN(query, label, values, -1)
}

func STextN[T any](query string, label string, values []T, n int) string {
	return strings.Replace(query, label,
		text.Join("'", "','", "','", "'", values),
		n,
	)
}

// AsText converts a byte array that may or may not be zlib-compressed into the original string
func AsText(b []byte) string {
	if len(b) > 1 && b[0] == 0x78 && b[1] == 0xda {
		r := bytes.NewReader(b)
		z, err := zlib.NewReader(r)
		if err == nil {
			defer z.Close()
			s, err := io.ReadAll(z)
			if err == nil {
				return string(s)
			}
		}
	}
	return string(b)
}
