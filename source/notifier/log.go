package notifier

import (
	"golib/lang"
	"golib/log"
)

type logNotifier struct {
	title string
}

func newLogNotifier(title string) *logNotifier {
	return &logNotifier{
		title: "[" + title + "]",
	}
}

func (self *logNotifier) Send(notificationType Type, title string, text string) (errors lang.Errors) {
	switch notificationType {
	case Info:
		log.Info(title, text)
	case Debug:
		log.Debug(title, text)
	case Warn:
		log.Warn(title, text)
	case Error:
		log.Error(title, text)
	}
	return
}
