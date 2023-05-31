package notifier

import (
	"github.com/ashep/sbk/telegram"
)

type Telegram struct {
	tg     *telegram.Client
	chatId string
}

func NewTelegram(tg *telegram.Client, chatId string) Notifier {
	return &Telegram{tg: tg, chatId: chatId}
}

func (n *Telegram) Notify(msg string) error {
	return n.tg.SendMessage(n.chatId, msg)
}
