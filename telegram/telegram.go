package telegram

type Client struct {
	Token string
}

func New(token string) *Client {
	return &Client{Token: token}
}
