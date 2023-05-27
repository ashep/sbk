package telegramnotifier

import (
	"github.com/ashep/sbk/telegram"
)

type Notifier struct {
	tg *telegram.Client
}

func New(tg *telegram.Client) *Notifier {
	return &Notifier{tg: tg}
}

func (n *Notifier) NotifySuccess(msg string) error {
	return nil
}

func (n *Notifier) NotifyError(msg string) error {
	return nil
}
