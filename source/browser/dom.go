package browser

import (
	"fmt"
	"strings"
	"time"

	"golib/arrays"
	"golib/lang"
	"golib/log"

	"github.com/tebeka/selenium"
)

type DOM interface {
	Find(by string, value string) *Node
	FindAll(by string, value string) []*Node
	FindIfPossible(by string, value string) *Node

	WaitFor(by string, value string) *Node
	WaitForEmpty(by string, value string) *Node
	WaitForOption(selectBy string, selectValue string, optionBy BySetting, optionMatch MatchSetting, optionValue string) *Node

	FindOption(selectBy string, selectValue string, optionBy BySetting, optionMatch MatchSetting, optionValue string) *Node
	FindOptionIfPossible(selectBy string, selectValue string, optionBy BySetting, optionMatch MatchSetting, optionValue string) *Node

	findWebElement(by string, value string) (selenium.WebElement, error)
	findWebElements(by string, value string) ([]selenium.WebElement, error)
	newNode(element selenium.WebElement) *Node
	browser() *Browser
}

// Browser (DOM)

func (self *Browser) Find(by string, value string) *Node {
	return Find(self, by, value)
}

func (self *Browser) FindAll(by string, value string) []*Node {
	return FindAll(self, by, value)
}

func (self *Browser) FindIfPossible(by string, value string) *Node {
	return FindIfPossible(self, by, value)
}

func (self *Browser) WaitForOneOf(description string, fns []func() bool) {
	WaitForOneOf(self, description, fns)
}

func (self *Browser) WaitFor(by string, value string) *Node {
	return WaitFor(self, by, value)
}

func (self *Browser) WaitForEmpty(by string, value string) *Node {
	return WaitForEmpty(self, by, value)
}

func (self *Browser) WaitForOption(selectBy string, selectValue string, optionBy BySetting, optionMatch MatchSetting, optionValue string) *Node {
	return WaitForOption(self, selectBy, selectValue, optionBy, optionMatch, optionValue)
}

func (self *Browser) FindOption(selectBy string, selectValue string, optionBy BySetting, optionMatch MatchSetting, optionValue string) *Node {
	return FindOption(self, selectBy, selectValue, optionBy, optionMatch, optionValue)
}

func (self *Browser) FindOptionIfPossible(selectBy string, selectValue string, optionBy BySetting, optionMatch MatchSetting, optionValue string) *Node {
	return FindOptionIfPossible(self, selectBy, selectValue, optionBy, optionMatch, optionValue)
}

func (self *Browser) findWebElement(by string, value string) (selenium.WebElement, error) {
	self.browser().DismissAlert()
	return self.wd.FindElement(by, value)
}

func (self *Browser) findWebElements(by string, value string) ([]selenium.WebElement, error) {
	self.browser().DismissAlert()
	return self.wd.FindElements(by, value)
}

func (self *Browser) newNode(element selenium.WebElement) *Node {
	return NewNode(self, element)
}

func (self *Browser) browser() *Browser {
	return self
}

// Node (DOM)

func (self *Node) Find(by string, value string) *Node {
	return Find(self, by, value)
}

func (self *Node) FindAll(by string, value string) []*Node {
	return FindAll(self, by, value)
}

func (self *Node) FindIfPossible(by string, value string) *Node {
	return FindIfPossible(self, by, value)
}

func (self *Node) WaitFor(by string, value string) *Node {
	return WaitFor(self, by, value)
}

func (self *Node) WaitForEmpty(by string, value string) *Node {
	return WaitForEmpty(self, by, value)
}

func (self *Node) WaitForOption(selectBy string, selectValue string, optionBy BySetting, optionMatch MatchSetting, optionValue string) *Node {
	return WaitForOption(self, selectBy, selectValue, optionBy, optionMatch, optionValue)
}

func (self *Node) FindOption(selectBy string, selectValue string, optionBy BySetting, optionMatch MatchSetting, optionValue string) *Node {
	return FindOption(self, selectBy, selectValue, optionBy, optionMatch, optionValue)
}

func (self *Node) FindOptionIfPossible(selectBy string, selectValue string, optionBy BySetting, optionMatch MatchSetting, optionValue string) *Node {
	return FindOptionIfPossible(self, selectBy, selectValue, optionBy, optionMatch, optionValue)
}

func (self *Node) findWebElement(by string, value string) (selenium.WebElement, error) {
	self.browser().DismissAlert()
	return self.element().FindElement(by, value)
}

func (self *Node) findWebElements(by string, value string) ([]selenium.WebElement, error) {
	self.browser().DismissAlert()
	return self.element().FindElements(by, value)
}

func (self *Node) newNode(element selenium.WebElement) *Node {
	return NewNode(self.browser(), element)
}

func (self *Node) browser() *Browser {
	return self._browser
}

// DOM

func Find(self DOM, by string, value string) *Node {
	log.Debug("Find " + by + " = [" + value + "]")
	element, err := self.findWebElement(by, value)
	assertNil(err)
	return self.newNode(element)
}

func FindAll(self DOM, by string, value string) []*Node {
	log.Debug("Find " + by + " = [" + value + "]")
	elements, err := self.findWebElements(by, value)
	assertNil(err)
	return arrays.Map(elements, func(element selenium.WebElement) *Node {
		return self.newNode(element)
	})
}

func FindIfPossible(self DOM, by string, value string) *Node {
	log.Debug("Find " + by + " = [" + value + "] (if possible)")
	element, err := self.findWebElement(by, value)
	if element == nil {
		log.Debug(strings.Split(err.Error(), "\n")[0])
		return nil
	}
	return self.newNode(element)
}

func WaitForOneOf(self DOM, description string, fns []func() bool) {
	WaitUntil(self, description, func() (time.Duration, bool) {
		for _, fn := range fns {
			if fn() {
				return 1, true
			}
		}
		return time.Second / 4, false
	})
}

func WaitFor(self DOM, by string, value string) *Node {
	var node *Node
	description := by + " = [" + value + "]"
	WaitUntil(self, description, func() (time.Duration, bool) {
		node = self.FindIfPossible(by, value)
		return time.Second / 4, node != nil && node.IsReady()
	})
	return node
}

func WaitForEmpty(self DOM, by string, value string) *Node {
	var node *Node
	description := "(EMPTY) " + by + " = [" + value + "]"
	WaitUntil(self, description, func() (time.Duration, bool) {
		node = self.FindIfPossible(by, value)
		return time.Second / 4, node != nil && node.IsReady() && node.IsEmpty()
	})
	return node
}

func WaitForOption(self DOM, selectBy string, selectValue string, optionBy BySetting, optionMatch MatchSetting, optionValue string) *Node {
	var node *Node
	description := selectBy + " = [" + selectValue + "] && " + string(optionBy) + " " + string(optionMatch) + " [" + optionValue + "]"
	log.Debug("Wait for " + description)
	WaitUntil(self, description, func() (time.Duration, bool) {
		node = self.FindOptionIfPossible(selectBy, selectValue, optionBy, optionMatch, optionValue)
		return time.Second / 4, node != nil && node.IsReady()
	})
	return node
}

func WaitUntil(self DOM, description string, durationGoalFunction func() (time.Duration, bool)) {
	for {
		success := lang.WaitUntilTimeout(durationGoalFunction, self.browser().WaitTimeout)
		if success {
			return
		}
		message := "Timeout: " + description
		b := self.browser()
		//b.sendWarning(message, lang.CurrentStack(0, 0))
		if b.CrashOnTimeout {
			log.Debug("Crashing due to timeout")
			panic(message)
		} else {
			wait := 4 * time.Minute
			log.Debug("Standing by for %s, retrying afterwards", wait)
			time.Sleep(wait)
		}
	}
}

func FindOption(self DOM, selectBy string, selectValue string, optionBy BySetting, optionMatch MatchSetting, optionValue string) *Node {
	var node *Node
	selectNode := self.Find(selectBy, selectValue)
	candidates, err := selectNode.element().FindElements(selenium.ByTagName, "option")
	assertNil(err)
	for _, candidate := range candidates {
		match, err := optionMatches(candidate, optionBy, optionMatch, optionValue)
		assertNil(err)
		if match {
			node = self.newNode(candidate)
			break
		}
	}
	assert(node != nil, "Option not found: "+string(optionBy)+" "+string(optionMatch)+" "+optionValue)
	return node
}

func FindOptionIfPossible(self DOM, selectBy string, selectValue string, optionBy BySetting, optionMatch MatchSetting, optionValue string) *Node {
	var node *Node
	selectNode := self.FindIfPossible(selectBy, selectValue)
	if selectNode != nil && selectNode.IsReady() {
		candidates, _ := selectNode.element().FindElements(selenium.ByTagName, "option")
		for _, candidate := range candidates {
			match, _ := optionMatches(candidate, optionBy, optionMatch, optionValue)
			if match {
				node = self.newNode(candidate)
				break
			}
		}
	}
	return node
}

type BySetting string

const (
	ByText  BySetting = "Text"
	ByValue BySetting = "Value"
)

type MatchSetting string

const (
	MatchStarts             MatchSetting = "Starts"
	MatchStartsIgnoreCase   MatchSetting = "StartsIgnoreCase"
	MatchEnds               MatchSetting = "Ends"
	MatchEndsIgnoreCase     MatchSetting = "EndsIgnoreCase"
	MatchContains           MatchSetting = "Contains"
	MatchContainsIgnoreCase MatchSetting = "ContainsIgnoreCase"
	MatchEquals             MatchSetting = "Equals"
	MatchEqualsIgnoreCase   MatchSetting = "EqualsIgnoreCase"
)

func optionMatches(option selenium.WebElement, optionBy BySetting, optionMatch MatchSetting, requiredContent string) (bool, error) {
	content, err := elementContent(option, optionBy)
	matches := err == nil && contentMatches(content, optionMatch, requiredContent)
	if !matches {
		log.Trace("Option %s [%s] doesn't %s [%s]", optionBy, content, optionMatch, requiredContent)
	}
	return matches, err
}

func contentMatches(content string, matchSetting MatchSetting, requiredContent string) bool {
	switch matchSetting {
	case MatchStarts:
		return strings.HasPrefix(content, requiredContent)
	case MatchStartsIgnoreCase:
		return lang.HasPrefixIgnoreCase(content, requiredContent)
	case MatchEnds:
		return strings.HasSuffix(content, requiredContent)
	case MatchEndsIgnoreCase:
		return lang.HasSuffixIgnoreCase(content, requiredContent)
	case MatchContains:
		return strings.Contains(content, requiredContent)
	case MatchContainsIgnoreCase:
		return lang.ContainsIgnoreCase(content, requiredContent)
	case MatchEquals:
		return content == requiredContent
	case MatchEqualsIgnoreCase:
		return strings.EqualFold(content, requiredContent)
	}
	panic(fmt.Sprintf("Unexpected MatchSetting: %s", matchSetting))
}

func elementContent(element selenium.WebElement, bySetting BySetting) (string, error) {
	switch bySetting {
	case ByText:
		text, err := element.Text()
		return text, err
	case ByValue:
		value, err := element.GetAttribute("value")
		return value, err
	}
	panic(fmt.Sprintf("Unexpected BySetting: %s", bySetting))
}
