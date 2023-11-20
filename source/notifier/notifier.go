package notifier

import "golib/lang"

type Notifier interface {
	Send(notificationType Type, title string, text string) lang.Errors
}

type Type int

const (
	Debug Type = iota
	Info
	Warn
	Error
)

func NewLog(prefix string) Notifier {
	return newLogNotifier(prefix)
}
