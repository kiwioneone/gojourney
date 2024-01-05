package gojourney

import (
	"context"
	"time"

	"github.com/kiwioneone/gojourney/discord"
)

func (c *Client) runJob(ctx context.Context, messageID string, customID string) (*CommandResult, error) {
	payload := NewMessageComponentPayload(c.GuildID, c.ChannelID, c.ws.SessionID(), messageID, discord.ComponentTypeButton, customID)

	messageFinder := &MessageFinder{Nonce: payload.Nonce}
	ob := NewCommonMessageObserver(messageFinder.FilterMessage)

	c.RegisterMessageObserver(ob)
	defer c.UnregisterMessageObserver(ob)

	err := c.sendHttpInteraction(ctx, payload)
	if err != nil {
		return nil, err
	}

	msg, err := ob.WaitMsg(time.Minute * 5)
	if err != nil {
		return nil, err
	}

	return MessageCommandResult(msg), err
}
