package gojourney

import (
	"context"

	"github.com/kiwioneone/gojourney/rest"

	"github.com/kiwioneone/gojourney/discord"
)

func (c *Client) GetMessage(ctx context.Context, messageID string) (*discord.Message, error) {
	if c.assistBot == nil {
		return nil, ErrorInvalidAssistBot
	}

	return rest.NewChannelHandler(c.assistBot).GetMessage(c.ChannelID, messageID)
}
