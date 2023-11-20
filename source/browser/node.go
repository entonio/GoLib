package browser

import (
	"strings"

	"golib/log"

	"github.com/tebeka/selenium"
)

type Node struct {
	_browser *Browser
	_element selenium.WebElement
}

func NewNode(browser *Browser, element selenium.WebElement) *Node {
	assert(element != nil, "WebElement cannot be nil")
	assert(browser != nil, "Browser cannot be nil")
	return &Node{
		_browser: browser,
		_element: element,
	}
}

func (self *Node) element() selenium.WebElement {
	return self._element
}

func (self *Node) String() string {
	result := self.Id()
	if len(result) == 0 {
		result = self.Name()
		if len(result) == 0 {
			value := self.Value()
			if len(value) > 0 {
				result = self.NodeName() + "(" + value + ")"
			} else {
				result = self.NodeName()
			}
		}
	}
	return result
}

func (self *Node) NodeName() string {
	s, err := self.element().TagName()
	if err != nil {
		log.Debug(err)
	}
	return s
}

func (self *Node) Id() string {
	return self.Attribute("id")
}

func (self *Node) Name() string {
	return self.Attribute("name")
}

func (self *Node) Title() string {
	return self.Attribute("title")
}

func (self *Node) HRef() string {
	return self.Attribute("href")
}

func (self *Node) Value() string {
	return self.Attribute("value")
}

func (self *Node) Attribute(name string) string {
	s, err := self.element().GetAttribute(name)
	if err != nil {
		log.Trace("%s(): %s", name, err)
	}
	return s
}

func (self *Node) Text() string {
	s, err := self.element().Text()
	if err != nil {
		log.Debug(err)
	}
	return s
}

func (self *Node) SetChecked(checked bool) *Node {
	if checked != self.IsChecked() {
		return self.Click()
	}
	return self
}

func (self *Node) Click() *Node {
	log.Debug("Click [%v]", self)
	err := self.element().Click()
	assertNil(err)
	return self
}

func (self *Node) Send(value string) *Node {
	log.Debug("SendKeys [%v] = [%v]", self, value)
	err := self.element().SendKeys(value)
	assertNil(err)
	return self
}

func (self *Node) SetText(value any) *Node {
	log.Debug("SetText [%v] = [%v]", self, value)
	err := self.element().Clear()
	assertNil(err)
	err = self.element().SendKeys(print(value))
	assertNil(err)
	return self
}

func (self *Node) IsChecked() bool {
	checked := self.Attribute("checked")
	//fmt.Println("SELECTED [" + selected + "]")
	return len(checked) > 0
}

func (self *Node) IsReady() bool {
	return self.IsDisplayed() && self.IsEnabled()
}

func (self *Node) IsDisplayed() bool {
	b, err := self.element().IsDisplayed()
	if err != nil {
		log.Debug(err)
	}
	return b
}

func (self *Node) IsEnabled() bool {
	b, err := self.element().IsEnabled()
	if err != nil {
		log.Debug(err)
	}
	return b
}

func (self *Node) IsEmpty() bool {
	text, err := self.element().Text()
	if err != nil {
		log.Debug(err)
	}
	return len(text) == 0
}

func (self *Node) Blur() {
	_, err := self.browser().wd.ExecuteScriptRaw(
		strings.Join([]string{
			"(function() {",
			"  var tmp = document.createElement('input');",
			"  document.body.appendChild(tmp);",
			"  tmp.focus();",
			"  document.body.removeChild(tmp);",
			"}) ();",
		}, "\n"),
		[]any{},
	)
	assertNil(err)
}
