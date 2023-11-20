package xmlbuilder

import (
	"fmt"
	"strings"
)

var FourSpaces = "    "
var Tab = "\t"

type Builder struct {
	buffer       []string
	stack        []string
	formatOutput bool

	Indent string
}

func New(formatOutput bool) *Builder {
	return WithEncoding("utf-8", formatOutput)
}

func WithEncoding(encoding string, formatOutput bool) *Builder {
	o := &Builder{formatOutput: formatOutput, Indent: FourSpaces}
	o.writeDeclaration(encoding)
	return o
}

func (self *Builder) writeDeclaration(encoding string) {
	self.push("<?xml version=\"1.0\" encoding=\"")
	self.push(encoding)
	self.push("\"?>")
	if self.formatOutput {
		self.push("\n")
	}
}

func (self *Builder) asStringBuilder() []string {
	if len(self.stack) > 0 {
		self.CloseAll()
	}
	return self.buffer
}

func (self *Builder) AsString() string {
	return strings.Join(self.asStringBuilder(), "")
}

func (self *Builder) push(s string) {
	self.buffer = append(self.buffer, s)
}

func (self *Builder) peek() string {
	return self.buffer[len(self.buffer)-1]
}

func (self *Builder) Open(nodeName string) *Builder {
	self.writeOpenTag(nodeName)
	self.stack = append(self.stack, nodeName)
	return self
}

func (self *Builder) Attributes(attributes ...any) *Builder {
	self.buffer = self.buffer[:len(self.buffer)-1]
	for i := 0; i < len(attributes); i += 2 {
		self.push(" ")
		self.push(fmt.Sprint(attributes[i]))
		self.push("=")
		self.push("\"")
		self.push(EscapeXml(fmt.Sprint(attributes[i+1])))
		self.push("\"")
	}
	self.push(">")
	return self
}

func (self *Builder) Close() *Builder {
	return self.closeLevels(1)
}

func (self *Builder) CloseAll() *Builder {
	return self.closeLevels(len(self.stack))
}

func (self *Builder) CloseAllExceptRoot() *Builder {
	levels := len(self.stack)
	if levels < 2 {
		return self
	}
	return self.closeLevels(levels - 1)
}

func (self *Builder) closeLevels(levels int) *Builder {
	limit := len(self.stack) - levels
	/*
		for i := len(self.stack) - 1; i >= limit; i -= 1 {
			nodeName := self.stack[i]
			self.writeCloseTag(nodeName)
		}
		for len(self.stack) > limit {
			self.stack = self.stack[:len(self.stack)-1]
		}
	*/
	for i := len(self.stack) - 1; i >= limit; i -= 1 {
		nodeName := self.stack[i]
		self.writeCloseTag(nodeName)
		self.stack = self.stack[:len(self.stack)-1]
	}
	return self
}

func (self *Builder) write(xml string) *Builder {
	self.push(xml)
	return self
}

func (self *Builder) PutIfTrue(nodeName string, value bool) *Builder {
	if value {
		return self.PutBool(nodeName, value)
	} else {
		return self
	}
}

func (self *Builder) PutIfExists(nodeName string, value any) *Builder {
	if value != nil {
		return self.Put(nodeName, value)
	} else {
		return self
	}
}

func (self *Builder) PutIfNotEmpty(nodeName string, value string) *Builder {
	if len(value) > 0 {
		return self.Put(nodeName, value)
	} else {
		return self
	}
}

func (self *Builder) PutBool(nodeName string, value bool) *Builder {
	self.writeOpenTag(nodeName)
	if value {
		self.push("true")
	} else {
		self.push("false")
	}
	self.writeCloseTag(nodeName)
	return self
}

func (self *Builder) PutInt(nodeName string, value int) *Builder {
	self.writeOpenTag(nodeName)
	self.push(fmt.Sprint(value))
	self.writeCloseTag(nodeName)
	return self
}

func (self *Builder) Put(nodeName string, value any) *Builder {
	s := fmt.Sprint(value)
	escaped := EscapeXml(s)
	return self.Literal(nodeName, escaped)
}

func (self *Builder) Text(value any) *Builder {
	s := fmt.Sprint(value)
	escaped := EscapeXml(s)
	self.push(escaped)
	return self
}

func (self *Builder) Literal(nodeName string, value string) *Builder {
	self.writeOpenTag(nodeName)
	self.push(value)
	self.writeCloseTag(nodeName)
	return self
}

func (self *Builder) writeOpenTag(nodeName string) {
	if self.formatOutput && self.peek() == ">" {
		self.push("\n")
		for i := 0; i < len(self.stack); i += 1 {
			self.push(self.Indent)
		}
	}
	self.push("<")
	self.push(nodeName)
	self.push(">")
}

func (self *Builder) writeCloseTag(nodeName string) {
	if nodeName == "Main" || nodeName == "Sub" {
		fmt.Printf("%s: %d\n", nodeName, len(self.stack))
	}
	if self.formatOutput && self.peek() == ">" {
		self.push("\n")
		for i := 0; i < len(self.stack)-1; i += 1 {
			self.push(self.Indent)
		}
	}
	self.push("</")
	self.push(nodeName)
	self.push(">")
}

func EscapeXml(original string) (s string) {
	for _, r := range original {

		// control and non-ascii
		if r < 0x20 || r > 0x7E {
			s += fmt.Sprintf("&#%d;", r)
			continue
		}

		c := fmt.Sprintf("%c", r)

		// xml special
		switch c {
		case "'":
			s += "&apos;"
		case "\"":
			s += "&quot;"
		case "<":
			s += "&lt;"
		case ">":
			s += "&gt;"
		case "&":
			s += "&amp;"
			// normal ascii
		default:
			s += c
		}
	}
	return
}
