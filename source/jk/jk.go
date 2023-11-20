package jk

import "encoding/json"

type Element struct {
	V any
}

var Format = json.MarshalIndent
var Print = json.Marshal

func Read(s string) *Element {
	m := map[string]any{}
	json.Unmarshal([]byte(s), &m)
	return &Element{V: m}
}

func (self *Element) K(keys ...string) *Element {
	v := self.V
	for _, key := range keys {
		v = v.(map[string]any)[key]
	}
	return &Element{V: v}
}

func (self *Element) I(indexes ...int) *Element {
	v := self.V
	for _, index := range indexes {
		v = v.([]any)[index]
	}
	return &Element{V: v}
}

func (self *Element) F(accept func(*Element) bool) *Element {
	for _, v := range self.V.([]any) {
		j := &Element{V: v}
		if accept(j) {
			return j
		}
	}
	return nil
}
