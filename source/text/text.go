package text

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"golib/lang"
)

func notFound(s string) error {
	return fmt.Errorf("[%s] not found", s)
}

func BeforeFirst(s string, delimiter string, inclusive bool) (string, error) {
	index := strings.Index(s, delimiter)
	if index < 0 {
		return "", notFound(delimiter)
	}
	if inclusive {
		return s[:index+len(delimiter)], nil
	} else {
		return s[:index], nil
	}
}

func BeforeLast(s string, delimiter string, inclusive bool) (string, error) {
	index := strings.LastIndex(s, delimiter)
	if index < 0 {
		return "", notFound(delimiter)
	}
	if inclusive {
		return s[:index+len(delimiter)], nil
	} else {
		return s[:index], nil
	}
}

func AfterFirst(s string, delimiter string, inclusive bool) (string, error) {
	index := strings.Index(s, delimiter)
	if index < 0 {
		return "", notFound(delimiter)
	}
	if inclusive {
		return s[index:], nil
	} else {
		return s[index+len(delimiter):], nil
	}
}

func AfterLast(s string, delimiter string, inclusive bool) (string, error) {
	index := strings.LastIndex(s, delimiter)
	if index < 0 {
		return "", notFound(delimiter)
	}
	if inclusive {
		return s[index:], nil
	} else {
		return s[index+len(delimiter):], nil
	}
}

func ConsonantalPrefix(word string) (prefix string) {
	for _, c := range word {
		if unicode.IsLetter(c) && !lang.IsVowel(c) {
			prefix += string(c)
		} else if len(prefix) > 0 {
			break
		}
	}
	return
}

func Syllables(word string) (syllables []string) {
	var syllable string
	var r rune

	afterVowel := false
	for _, r = range word {
		if lang.IsVowel(r) {
			afterVowel = true
		} else {
			if afterVowel {
				syllables = append(syllables, syllable)
				syllable = ""
			}
			afterVowel = false
		}
		syllable += string(r)
	}

	previous := len(syllables) - 1
	if previous < 0 || lang.IsVowel(r) {
		syllables = append(syllables, syllable)
	} else {
		syllables[previous] = syllables[previous] + syllable
	}

	return
}

func ReplaceAllStringSubmatchFunc(re *regexp.Regexp, source string, replaceFn func(match []string) string) string {
	result := ""
	lastIndex := 0

	for _, v := range re.FindAllSubmatchIndex([]byte(source), -1) {
		groups := []string{}
		for i := 0; i < len(v); i += 2 {
			if v[i] == -1 || v[i+1] == -1 {
				groups = append(groups, "")
			} else {
				groups = append(groups, source[v[i]:v[i+1]])
			}
		}

		result += source[lastIndex:v[0]] + replaceFn(groups)
		lastIndex = v[1]
	}

	return result + source[lastIndex:]
}

func Join(prefix string, separator string, conjunction string, suffix string, arrayOrObject any) string {
	elementos := lang.Strings(arrayOrObject)
	text := prefix
	limit := len(elementos) - 2
	for i, s := range elementos {
		text += fmt.Sprint(s)
		if i < limit {
			text += separator
		} else if i == limit {
			text += conjunction
		}
	}
	return text + suffix
}
