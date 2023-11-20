package lang

import (
	"bytes"
	"io"
	"reflect"
	"regexp"
	"strings"
	"unicode"

	encoding_unicode "golang.org/x/text/encoding/unicode"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func AsNormalisedKey(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "/", " ")
	s = strings.ReplaceAll(s, "(", " ")
	s = strings.ReplaceAll(s, ")", " ")
	s = strings.ReplaceAll(s, "-", " ")
	s = strings.ReplaceAll(s, ",", " ")
	s = strings.ReplaceAll(s, "  ", " ")
	s = strings.ReplaceAll(s, "  ", " ")
	s = strings.ReplaceAll(s, "s ", " ")
	s = strings.TrimSuffix(s, "s")
	s = strings.ReplaceAll(s, "ãe", "ao")
	s = strings.ReplaceAll(s, "ão", "ao")
	s = strings.ReplaceAll(s, "õe", "ao")
	s = RemoveDiacritics(s)
	s = strings.ReplaceAll(s, " a ", " ")
	s = strings.ReplaceAll(s, " e ", " ")
	s = strings.ReplaceAll(s, " o ", " ")
	s = strings.ReplaceAll(s, " da ", " ")
	s = strings.ReplaceAll(s, " de ", " ")
	s = strings.ReplaceAll(s, " do ", " ")
	return s
}

func Mask(s string, values ...string) string {
	for _, k := range values {
		s = strings.ReplaceAll(s, k, "*******")
	}
	return s
}

func ReplacePrefix(s string, prefix string, replacement string) string {
	if strings.HasPrefix(s, prefix) {
		return replacement + s[len(prefix):]
	} else {
		return s
	}
}

func ReplaceSuffix(s string, suffix string, replacement string) string {
	if strings.HasSuffix(s, suffix) {
		return s[0:len(s)-len(suffix)] + replacement
	} else {
		return s
	}
}

func RotateString(s string, key int) string {
	result := ""
	for _, r := range s {
		result += string(rotateRune(r, key))
	}
	return result
}

func rotateRune(r rune, key int) rune {
	if r >= 'A' && r <= 'Z' {
		base := int('A')
		return rune((int(r)-base-2+key)%26 + base)
	}
	if r >= 'a' && r <= 'z' {
		base := int('a')
		return rune((int(r)-base-2+key)%26 + base)
	}
	return r
}

func flips(input string) (output string) {
	var out []rune
	for _, c := range input {
		out = append(out, flipc(c))
	}
	output = string(out)
	return
}

func flipc(input rune) rune {
	switch {
	case input > 50 && input < 55:
		return 105 - input
	case input > 70 && input < 89:
		return 159 - input
	case input > 100 && input < 121:
		return 221 - input
	case input == 48:
		return 56
	case input == 56:
		return 48
	}
	return input
}

func SplitFixed(original string, width int) (split []string) {
	first := len(original) % width
	parts := len(original) / width
	if first > 0 {
		split = append(split, original[0:first])
	}
	for i := 0; i < parts; i++ {
		next := first + width
		split = append(split, original[first:next])
		first = next
	}
	return
}

func AsEmoji(b bool) string {
	return As(b, "✅", "❌").(string)
}

func Strings(value any) []string {
	switch value.(type) {
	case []string:
		return value.([]string)
	case string:
		return []string{value.(string)}
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.Array, reflect.Slice:
		output := make([]string, v.Len())
		for i := 0; i < len(output); i += 1 {
			output[i] = Print(v.Index(i))
		}
		return output
	}
	return []string{Print(value)}
}

var tRemoveDiacritics = transform.Chain(
	norm.NFD,
	runes.Remove(runes.In(unicode.Mn)), // Mn is 'nonspacing marks'
	norm.NFC,
)

func RemoveDiacritics(s string) string {
	t, _, _ := transform.String(tRemoveDiacritics, s)
	return t
}

func KeyFromString(s string) string {
	s = strings.ReplaceAll(s, "*", "")
	s = strings.TrimSpace(s)
	s = strings.Trim(s, ".?!:;")
	s = strings.TrimSpace(s)
	s = strings.ToUpper(s)
	s = RemoveDiacritics(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "--", "-")
	s = strings.ReplaceAll(s, "-DE-", "-")
	s = strings.ReplaceAll(s, "-DA-", "-")
	s = strings.ReplaceAll(s, "-DAS-", "-")
	s = strings.ReplaceAll(s, "-DO-", "-")
	s = strings.ReplaceAll(s, "-DOS-", "-")
	return s
}

func StringFromWindowsBytes(source []byte) (string, error) {
	// Make a transformer that converts MS-Win default to UTF8
	win16le := encoding_unicode.UTF16(encoding_unicode.LittleEndian, encoding_unicode.IgnoreBOM)
	// Make a transformer that is like win16le, but abides by BOM
	utf16bom := encoding_unicode.BOMOverride(win16le.NewDecoder())
	// Make a Reader that uses utf16bom
	reader := transform.NewReader(bytes.NewReader(source), utf16bom)
	// decode
	decoded, err := io.ReadAll(reader)
	return string(decoded), err
}

func HasPrefixIgnoreCase(s, prefix string) bool {
	return len(s) >= len(prefix) && strings.EqualFold(s[0:len(prefix)], prefix)
}

func HasSuffixIgnoreCase(s, suffix string) bool {
	return len(s) >= len(suffix) && strings.EqualFold(s[len(s)-len(suffix):], suffix)
}

func ContainsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) && strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// Instead of considering the []byte as the UTF-8 representation of the string,
// convert each byte into the character it represents.
// Other 8-bit encodings might be handled by doing character translation.
func StringFromIsoLatin1Bytes(bytes []byte) string {
	runes := make([]rune, len(bytes))
	for i, iso_8859_1_byte := range bytes {
		runes[i] = rune(iso_8859_1_byte)
	}
	return string(runes)
}

func Slice(s string, re *regexp.Regexp, n int) []string {
	if n == 0 {
		return nil
	}

	matches := re.FindAllStringIndex(s, n)
	strings := make([]string, 0, len(matches))

	beg := 0
	end := 0
	for _, match := range matches {
		if n > 0 && len(strings) >= n-1 {
			break
		}

		end = match[0]
		if match[1] != 0 {
			strings = append(strings, s[beg:end])
		}
		beg = match[1]
		// This also appends the current match
		strings = append(strings, s[match[0]:match[1]])
	}

	if end != len(s) {
		strings = append(strings, s[beg:])
	}

	return strings
}

func TrimTo(s string, max int, ellipsis string) string {
	if len(s) <= max {
		return s
	}
	return s[0:max-len(ellipsis)] + ellipsis
}

func IsVowel(r rune) bool {
	b := RemoveDiacritics(string(r))[0]
	switch b {
	case 'a', 'e', 'i', 'o', 'u', 'w', 'y',
		'A', 'E', 'I', 'O', 'U', 'W', 'Y':
		return true
	default:
		return false
	}
}

func StringFromBool(b bool, ifTrue string, ifFalse string) string {
	if b {
		return ifTrue
	} else {
		return ifFalse
	}
}
