package telegram

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	Token string

	cli *http.Client
}

func New(token string) *Client {
	return &Client{
		Token: token,
		cli:   &http.Client{},
	}
}

func (c *Client) SendMessage(chatId, msg string) error {
	payload := fmt.Sprintf(`{"chat_id":%q,"text":%q,"parse_mode":"markdown"}`, chatId, msg)

	fmt.Printf("%s\n", payload)

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", c.Token),
		strings.NewReader(payload),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := c.cli.Do(req)
	if err != nil {
		return err
	}

	rb, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("%s", rb)
	}

	return nil
}
