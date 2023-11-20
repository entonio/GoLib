package cs

import "strings"

func Min(i1 int, i2 int) int {
	if i1 < i2 {
		return i1
	} else {
		return i2
	}
}

func Max(i1 int, i2 int) int {
	if i1 > i2 {
		return i1
	} else {
		return i2
	}
}

func Length(s *string) *int {
	if s == nil {
		return nil
	}
	t := len(*s)
	return &t
}

func IndexOf(s *string, v string) *int {
	if s == nil {
		return nil
	}
	t := strings.Index(*s, v)
	return &t
}

func Replace(s *string, v1 string, v2 string) *string {
	if s == nil {
		return nil
	}
	t := strings.ReplaceAll(*s, v1, v2)
	return &t
}

func ReplaceSuffix(s *string, v1 string, v2 string) *string {
	if s == nil {
		return nil
	}
	t := *s
	if strings.HasSuffix(t, v1) {
		t = t[:len(t)-len(v1)] + v2
	}
	return &t
}

func Substring(s *string, start int, length int) *string {
	if s == nil {
		return nil
	}
	t := (*s)[start:length]
	return &t
}

func ToLower(s *string) *string {
	if s == nil {
		return nil
	}
	t := strings.ToLower(*s)
	return &t
}

func Trim(s *string) *string {
	if s == nil {
		return nil
	}
	t := strings.TrimSpace(*s)
	return &t
}
