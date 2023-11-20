package settings

import (
	"encoding/base64"
	"regexp"
)

type Settings interface {
	B(section string, key string) bool
	I(section string, key string) int
	F(section string, key string) float64
	S(section string, key string) string
	SL(section string, key string) []string
	C(section string, key string) string
}

func NewIni(path string) Settings {
	return newIni(path)
}

func SL(csv string) []string {
	return regexp.MustCompile(`\s*,\s*`).Split(csv, -1)
}

func C(b64 string) string {
	s, err := base64.StdEncoding.DecodeString(b64)
	assertNil(err)
	return string(s)
}
