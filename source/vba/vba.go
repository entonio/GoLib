package vba

import (
	"fmt"
	"strconv"

	"golib/cs"
)

func INSTR(s1 any, s2 string) *int {
	var zero = 0
	var start = 1
	var result *int
	if s1 == nil {
		result = nil
	} else if *_len(s1) == 0 {
		result = &zero
	} else if *_len(s2) == 0 {
		result = &start
	} else {
		var position = *cs.IndexOf(Str(s1), s2) + 1
		if position < start {
			result = &zero
		} else {
			result = &position
		}
	}
	fmt.Printf("INSTR([%s]; [%s]) -> [%s]", s1, s2, result)
	return result
}

func LCASE(val any) *string {
	return cs.ToLower(Str(val))
}

func LEN(val any) *int {
	result := _len(val)
	fmt.Printf("LEN(\"%s\") -> %d", val, result)
	return result
}

func _len(val any) *int {
	return cs.Length(Str(val))
}

func LEFT(val any, length int) *string {
	var result *string
	var s = Str(val)
	if s == nil {
		result = nil
	} else {
		var len = cs.Min(length, *_len(val))
		result = cs.Substring(s, 0, len)
	}
	fmt.Printf("LEFT([%s]; %s) -> [%s]", val, length, result)
	return result
}

func RIGHT(val any, length int) *string {
	var result *string
	var s = Str(val)
	if s == nil {
		result = nil
	} else {
		var max = *_len(val)
		length = cs.Min(max, length)
		var start = 1 + cs.Max(0, max-length)
		//fmt.Printf("Substring([{val}]; {start - 1}; {length})");
		result = cs.Substring(s, start-1, length)
	}
	fmt.Printf("RIGHT([%s]; %d) -> [%s]", val, length, *result)
	return result
}

func MID(val any, start int, length *int) *string {
	result := _mid(val, start, length)
	fmt.Printf("MID([%s]; %d; %d) -> [%s]", val, start, length, result)
	return result
}

func _mid(val any, start int, length *int) *string {
	var result *string
	var s = Str(val)
	if s == nil {
		result = nil
	} else {
		var max = *_len(val)
		if length == nil {
			length = &max
		}
		min := cs.Min(*length, max-start+1)
		result = cs.Substring(s, start-1, min)
	}
	return result
}

func REPLACE(val any, v1 string, v2 string) *string {
	return cs.Replace(Str(val), v1, v2)
}

func TRIM(val any) *string {
	return cs.Trim(Str(val))
}

func Str(val any) *string {
	if val == nil {
		return nil
	}
	s := fmt.Sprint(val)
	return &s
}

func Int(val any) *int {
	var result *int
	var s = *cs.Trim(Str(val))
	if "" == s {
		result = nil
	} else {
		i, err := strconv.Atoi(s)
		if err != nil {
			result = nil
		} else {
			result = &i
		}
	}
	fmt.Printf("Int([%s]) -> %d", val, result)
	return result
}

func Double(val any) *float64 {
	var result *float64
	var s = *cs.Trim(Str(val))
	if "" == s {
		result = nil
	} else {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			result = nil
		} else {
			result = &f
		}
	}
	fmt.Printf("Double([%s]) -> %f", val, result)
	return result
}
