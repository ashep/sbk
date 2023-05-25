package telegram

type Telegram struct {
	Token string
}

func New(token string) *Telegram {
	return &Telegram{Token: token}
}
