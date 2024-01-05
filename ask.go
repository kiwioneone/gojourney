package gojourney

import (
	"context"
	"time"
)

func NewAskCommandPayload(guildID, channelID, sessionID string, question string) *Payload {
	return NewApplicationCommandPayload(
		guildID,
		channelID,
		sessionID,
		"ask",
		"question",
		question,
		ApplicationCommandIdVersion{
			Version: "1166847114203123794",
			ID:      "994261739745050684",
		})
}

func (c *Client) Ask(ctx context.Context, question string) (result string, err error) {

	payload := NewAskCommandPayload(c.GuildID, c.ChannelID, c.ws.SessionID(), question)
	messageFinder := &MessageFinder{Nonce: payload.Nonce, NoJob: true}

	ob := NewCommonMessageObserver(messageFinder.FilterMessage)
	c.RegisterMessageObserver(ob)
	defer c.UnregisterMessageObserver(ob)

	err = c.sendHttpInteraction(ctx, payload)
	if err != nil {
		return "", err
	}

	msg, err := ob.WaitMsg(time.Minute * 5)
	if err != nil {
		return "", err
	}

	if len(msg.Embeds) > 0 {
		return msg.Embeds[0].Description, nil
	}

	return "", nil
}
