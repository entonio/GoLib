package lang

import (
	"fmt"

	"golib/jk"
)

func Print(s any) string {
	switch s.(type) {
	case string:
		return s.(string)
	default:
		return fmt.Sprint(s)
	}
}

func ToString[A any](a A) string {
	return Print(a)
}

func Pretty(value any) string {
	bytes, err := jk.Format(value, "", "\t")
	if bytes != nil {
		return string(bytes)
	}
	return err.Error()
}
