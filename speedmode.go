package gojourney

import (
	"context"
	"errors"
	"time"
)

func NewFastCommandPayload(guildID, channelID, sessionID string) *Payload {
	return NewApplicationCommandPayload(
		guildID,
		channelID,
		sessionID,
		"fast",
		"",
		"",
		ApplicationCommandIdVersion{
			Version: "987795926183731231",
			ID:      "972289487818334212",
		})
}

func (c *Client) Fast(ctx context.Context) error {
	payload := NewFastCommandPayload(c.GuildID, c.ChannelID, c.ws.SessionID())
	messageFinder := &MessageFinder{Nonce: payload.Nonce, NoJob: true}

	ob := NewCommonMessageObserver(messageFinder.FilterMessage)
	c.RegisterMessageObserver(ob)
	defer c.UnregisterMessageObserver(ob)

	err := c.sendHttpInteraction(ctx, payload)
	if err != nil {
		return err
	}

	msg, err := ob.WaitMsg(time.Minute * 5)
	if err != nil {
		return err
	}

	if MessageEmbedDescription(msg) != "" {
		return errors.New(MessageEmbedDescription(msg))
	}

	return nil
}

func NewRelaxCommandPayload(guildID, channelID, sessionID string) *Payload {
	return NewApplicationCommandPayload(
		guildID,
		channelID,
		sessionID,
		"relax",
		"",
		"",
		ApplicationCommandIdVersion{
			Version: "987795926183731232",
			ID:      "972289487818334213",
		})
}

func (c *Client) Relax(ctx context.Context) error {
	payload := NewRelaxCommandPayload(c.GuildID, c.ChannelID, c.ws.SessionID())
	messageFinder := &MessageFinder{Nonce: payload.Nonce, NoJob: true}

	ob := NewCommonMessageObserver(messageFinder.FilterMessage)
	c.RegisterMessageObserver(ob)
	defer c.UnregisterMessageObserver(ob)

	err := c.sendHttpInteraction(ctx, payload)
	if err != nil {
		return err
	}

	msg, err := ob.WaitMsg(time.Minute * 5)
	if err != nil {
		return err
	}

	if MessageEmbedDescription(msg) != "" {
		return errors.New(MessageEmbedDescription(msg))
	}

	return nil
}
