package browser

import (
	"fmt"
	"strings"
	"time"

	"golib/log"
	"golib/notifier"

	"github.com/tebeka/selenium"
)

type Engine string

const (
	CHROME  Engine = "chrome"
	FIREFOX Engine = "firefox"
)

type Option string

const (
	Op_AcceptInsecureCerts Option = "Op_AcceptInsecureCerts"
	Op_AutoClose           Option = "Op_AutoClose"
	Op_NoClientCert        Option = "Op_NoClientCert"
	Op_NoWindow            Option = "Op_NoWindow"
	Op_SilentLaunch        Option = "Op_SilentLaunch"
	Op_ProcessPerSite      Option = "Op_ProcessPerSite"
	Op_ProcessPerTab       Option = "Op_ProcessPerTab"
	Op_SingleProcess       Option = "Op_SingleProcess"
	Op_Offscreen           Option = "Op_Offscreen"
)

type Browser struct {
	wd selenium.WebDriver

	UrlTransforms  []UrlTransform
	Notifier       notifier.Notifier
	WaitTimeout    time.Duration
	CrashOnTimeout bool
}

type UrlTransform func(url string) string

const (
	secondsAfterSwitchWindow = 2
	secondsAfterSwitchFrame  = 3.5
)

func New(wd selenium.WebDriver) *Browser {
	assert(wd != nil, "WebDriver cannot be nil")
	return &Browser{
		wd:          wd,
		Notifier:    notifier.NewLog("BROWSER NOTIFICATION"),
		WaitTimeout: time.Second * 120,
	}
}

func (self *Browser) Close() error {
	return self.wd.Quit()
}

func (self *Browser) Wait(seconds float64) {
	time.Sleep(time.Duration(seconds * float64(time.Second)))
}

func (self *Browser) Title(title string) {
	_, err := self.wd.ExecuteScriptRaw(
		fmt.Sprintf("document.title='%s'", strings.ReplaceAll(title, "'", "\\'")),
		[]any{},
	)
	assertNil(err)
}

func (self *Browser) CreateTab() {
	_, err := self.wd.ExecuteScriptRaw(
		"window.open()",
		[]any{},
	)
	assertNil(err)
}

func (self *Browser) CurrentTabHandle() string {
	handle, err := self.wd.CurrentWindowHandle()
	assertNil(err)
	log.Debug("Current tab is [%s]", handle)
	return handle
}

func (self *Browser) TabWithHandle(handle string) {
	log.Debug("Will switch to tab [%s]", handle)
	err := self.wd.SwitchWindow(handle)
	assertNil(err)
	log.Debug("Waiting " + print(secondsAfterSwitchWindow) + " seconds for DOM to update")
	self.Wait(secondsAfterSwitchWindow)
}

func (self *Browser) FirstTab() {
	self.TabNumber(0)
}

func (self *Browser) LatestTab() {
	self.TabNumber(-1)
}

func (self *Browser) TabNumber(index int) {
	handles, err := self.wd.WindowHandles()
	assertNil(err)
	if index < 0 {
		index = len(handles) + index
	}
	self.TabWithHandle(handles[index])
}

func (self *Browser) FrameTop() {
	err := self.wd.SwitchFrame(nil)
	assertNil(err)
}

func (self *Browser) FrameIndex(indices ...int) {
	assert(len(indices) > 0, "At least one frame index required")
	for _, index := range indices {
		log.Debug("Will switch to frame " + print(index))
		err := self.wd.SwitchFrame(index)
		assertNil(err)
		log.Debug("Waiting " + print(secondsAfterSwitchFrame) + " seconds for DOM to update")
		self.Wait(secondsAfterSwitchFrame)
	}
}

func (self *Browser) FrameName(names ...string) {
	assert(len(names) > 0, "At least one frame name required")
	for _, name := range names {
		/*
			log.Debug("Will switch to frame [" + name + "]")
			err := self.wd.SwitchFrame(name)
		*/
		log.Debug("Will locate frame [" + name + "]")
		frame, err := self.findWebElement(selenium.ByCSSSelector, "frame[name="+name+"], iframe[name="+name+"]")
		assertNil(err)
		tag, err := frame.TagName()
		assertNil(err)
		log.Debug("Will switch to " + tag)
		err = self.wd.SwitchFrame(frame)
		assertNil(err)
		log.Debug("Waiting " + print(secondsAfterSwitchFrame) + " seconds for DOM to update")
		self.Wait(secondsAfterSwitchFrame)
	}
}

func (self *Browser) Get(url string, parameters ...any) {
	for i := 0; i < len(parameters); i += 2 {
		url = strings.ReplaceAll(url, print(parameters[i]), print(parameters[i+1]))
	}
	for _, transform := range self.UrlTransforms {
		url = transform(url)
	}
	err := self.wd.Get(url)
	assertNil(err)
}

func (self *Browser) DismissAlert() {
	text, _ := self.wd.AlertText()
	if len(text) > 0 {
		log.Debug("DISMISSED ALERT: [" + text + "]")
	}
	self.wd.DismissAlert()
}

func (self *Browser) AcceptAlert() {
	text, _ := self.wd.AlertText()
	if len(text) > 0 {
		log.Debug("DISMISSED ALERT: [" + text + "]")
	}
	self.wd.AcceptAlert()
}

/*
func (self *Browser) sendWarning(title string, message string) {
	self.Notifier.Send(notifier.Warn, title, message)
}
*/

func XPathTextContains(text string) string {
	return "//*[contains(text(),'" + text + "')]"
}

func XPathTextContainsIgnoreCase(text string) string {
	//	return "(//text()[contains(translate(., 'ABCDEFGHIJKLMNOPQRSTUVWXYZ', 'abcdefghijklmnopqrstuvwxyz'), '" + strings.ToLower(text) + "')])[1]/.."
	return "(//text()[contains(translate(., 'ABCDEFGHIJKLMNOPQRSTUVWXYZÁÀÂÃÉÈÊÍÌÓÒÔÕÚÙÇ', 'abcdefghijklmnopqrstuvwxyzáàâãéèêíìóòôõúùç'), '" + strings.ToLower(text) + "')])/.."
}
