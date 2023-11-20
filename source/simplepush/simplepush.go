package simplepush

import (
	"golib/lang"
	"golib/notifier"
	libsettings "golib/settings"
)

type SimplePush struct {
	active      bool
	AppKey      string
	UserKey     string
	InfoDevices []string
	WarnDevices []string
}

func NewSimplePush() *SimplePush {
	settings := libsettings.NewIni("configs/settings.ini")
	sp := SimplePush{
		AppKey:      settings.S("SimplePush", "AppKey"),
		UserKey:     settings.S("SimplePush", "UserKey"),
		InfoDevices: settings.SL("SimplePush", "InfoDevices"),
		WarnDevices: settings.SL("SimplePush", "WarnDevices"),
	}
	sp.active = (len(sp.InfoDevices) > 0 || len(sp.WarnDevices) > 0) && len(sp.AppKey) > 0 && len(sp.UserKey) > 0
	return &sp
}

func (self *SimplePush) Send(notificationType notifier.Type, label string, title string, message string) (errors lang.Errors) {
	if !self.active {
		return
	}
	var devices []string
	if notificationType <= notifier.Info {
		devices = append(devices, self.InfoDevices...)
	} else if notificationType <= notifier.Error {
		devices = append(devices, self.WarnDevices...)
	}
	return notifier.NewPushover(label+" ", self.AppKey, self.UserKey, devices).
		Send(notificationType, title, message)
}
