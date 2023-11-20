package notifier

import (
	"fmt"

	"golib/lang"
	"golib/log"

	"github.com/gregdel/pushover"
)

func NewPushover(label string, appKey string, userKey string, devices []string) Notifier {
	return newPushoverNotifier(label, appKey, userKey, devices)
}

type pushoverNotifier struct {
	label   string
	appKey  string
	userKey string
	devices []string
}

var soundForType = map[Type]string{
	Debug: pushover.SoundPianobar,
	Info:  pushover.SoundMagic,
	Warn:  pushover.SoundGamelan,
	Error: pushover.SoundSpaceAlarm,
}

func newPushoverNotifier(label string, appKey string, userKey string, devices []string) *pushoverNotifier {
	return &pushoverNotifier{
		label:   label,
		appKey:  appKey,
		userKey: userKey,
		devices: devices,
	}
}

func (self *pushoverNotifier) Send(notificationType Type, title string, text string) (errors lang.Errors) {
	for _, device := range self.devices {
		e := self.send(notificationType, title, text, device, true)
		for _, v := range e.Map {
			errors.Add(v)
		}
	}
	return errors
}

func (self *pushoverNotifier) send(notificationType Type, title string, text string, device string, warn bool) (errors lang.Errors) {
	log.Debug("Sending push notification to %s, title=[%s], message=[%s]...", device, title, text)
	// strings.ReplaceAll(text, "\n", " "))

	max := 512
	if len(self.label) > 0 {
		title = self.label + title
	}
	title = lang.TrimTo(title, 180, "…")
	text = lang.TrimTo(text, max-len(title), "…")

	// Send the message to the recipient
	app := pushover.New(self.appKey)
	recipient := pushover.NewRecipient(self.userKey)
	message := pushover.NewMessageWithTitle(text, title)
	message.DeviceName = device //strings.Join(self.devices, ",")
	message.Sound = soundForType[notificationType]
	response, err := app.SendMessage(message, recipient)
	if err != nil {
		log.Error("Error sending message: %v\nTitle:[%s]\nText:[%s]", err, title, text)
		errors.Add(err)
		return
	}

	for _, e := range response.Errors {
		log.Error(err)
		errors.Add(errors.AsError(e))
	}

	log.Debug("Used %d of %d push notifications", response.Limit.Total-response.Limit.Remaining, response.Limit.Total)
	if warn && (response.Limit.Remaining == 100 || response.Limit.Remaining == 50) {
		self.Send(
			Warn,
			"Message limit is nigh!",
			fmt.Sprintf("There are less than %d messages left until %s",
				response.Limit.Remaining,
				response.Limit.NextReset,
			),
		)
	}

	return
}
