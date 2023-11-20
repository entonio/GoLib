package lang

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func As(b bool, onTrue any, onFalse any) any {
	if b {
		return onTrue
	} else {
		return onFalse
	}
}

func If[V any](condition bool, ifTrue V, ifFalse V) V {
	if condition {
		return ifTrue
	} else {
		return ifFalse
	}
}

func AsBool(s string) bool {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "t", "true", "y", "yes", "1":
		return true
	case "f", "false", "n", "no", "0":
		return false
	}
	panic("Unable to convert <" + s + "> to bool")
}

var regexpNumeric = regexp.MustCompile(`^\d+(\.[\d]+)?$`)

var regexpAllExceptDigits = regexp.MustCompile(`[^\d]+`)

func Digits(s string) string {
	return regexpAllExceptDigits.ReplaceAllString(s, "")
}

func IsNumeric(s string) bool {
	return regexpNumeric.MatchString(s)
}

func AsInt(s string) (i int) {
	i, err := strconv.Atoi(s)
	AssertNil(err)
	return
}

func ForceInt(s string, alternative int) int {
	i, err := strconv.Atoi(s)
	if err == nil {
		return i
	} else {
		return alternative
	}
}

func ConstrainInt(i int, min int, max int) int {
	if i < min {
		return min
	} else if i > max {
		return max
	} else {
		return i
	}
}

func MinInt(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func MaxInt(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func SameFloat(a float64, b float64) bool {
	d := Abs(a - b)
	return d < 0.000001 // float32 doesn't ensure more than 6 significant digits
}

func Abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	if x == 0 {
		return 0
	}
	return x
}

func PrintFloat(value float64, minPrecision int, maxPrecision int, thousands string, decimals string) string {
	var s string
	if maxPrecision >= 0 {
		s = fmt.Sprintf("%."+Print(maxPrecision)+"f", value)
	} else {
		s = fmt.Sprintf("%f", value)
	}

	split := strings.Split(s, ".")

	integers := SplitFixed(split[0], 3)

	var fractional string
	if len(split) == 2 {
		fractional = strings.TrimRight(split[1], "0")
	}

	for len(fractional) < minPrecision {
		fractional += "0"
	}

	if len(fractional) > 0 {
		return strings.Join(integers, thousands) + decimals + fractional
	} else {
		return strings.Join(integers, thousands)
	}
}
