package gojourney

import "github.com/kiwioneone/gojourney/rest"

type Option func(client *Client)

func WithAssistBot(botToken string) Option {
	return func(client *Client) {
		client.assistBot = rest.NewClient(botToken)
	}
}
