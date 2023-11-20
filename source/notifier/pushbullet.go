//go:build discontinued

package notifier

import (
	"fmt"

	"golib/lang"
	"golib/log"

	pushbullet "github.com/xconstruct/go-pushbullet"
)

func NewPushbullet(label string, tokens []string) Notifier {
	return newPushbulletNotifier(label, tokens)
}

type pushbulletNotifier struct {
	label  string
	tokens []string
}

func newPushbulletNotifier(label string, tokens []string) *pushbulletNotifier {
	return &pushbulletNotifier{
		label:  label,
		tokens: tokens,
	}
}

func (self *pushbulletNotifier) Send(notificationType Type, title string, text string) (errors lang.Errors) {
	for _, token := range self.tokens {
		errors.Set(token, self.sendTo(token, title, text))
	}
	return
}

func (self *pushbulletNotifier) sendTo(token string, title string, text string) error {
	pb := pushbullet.New(token)

	devices, err := pb.Devices()
	if err != nil {
		log.Error(err)
		return err
	}

	if len(devices) == 0 {
		log.Debug("Found no device")
		return fmt.Errorf("Found no device")
	}

	if len(self.label) > 0 {
		title = self.label + title
	}

	err = pb.PushNote(devices[0].Iden, title, text)
	if err != nil {
		log.Debug(err)
		return err
	}

	return nil
}
