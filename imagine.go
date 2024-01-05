package gojourney

import (
	"context"
	"time"
)

func NewImagineCommandPayload(guildID, channelID, sessionID string, prompt string) *Payload {
	return NewApplicationCommandPayload(
		guildID,
		channelID,
		sessionID,
		"imagine",
		"prompt",
		prompt,
		ApplicationCommandIdVersion{
			Version: "1166847114203123795",
			ID:      "938956540159881230",
		})
}

func (c *Client) Imagine(ctx context.Context, prompt string) (result *CommandResult, err error) {

	payload := NewImagineCommandPayload(c.GuildID, c.ChannelID, c.ws.SessionID(), prompt)
	messageFinder := &MessageFinder{Nonce: payload.Nonce}

	ob := NewCommonMessageObserver(messageFinder.FilterMessage)
	c.RegisterMessageObserver(ob)
	defer c.UnregisterMessageObserver(ob)

	err = c.sendHttpInteraction(ctx, payload)
	if err != nil {
		return nil, err
	}

	msg, err := ob.WaitMsg(time.Minute * 5)
	if err != nil {
		return nil, err
	}

	return MessageCommandResult(msg), nil
}
