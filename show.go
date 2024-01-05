package gojourney

import (
	"context"
	"time"

	"github.com/kiwioneone/gojourney/discord"
)

func NewShowCommandPayload(guildID, channelID, sessionID string, jobID string) *Payload {
	return NewApplicationCommandPayload(
		guildID, channelID, sessionID,
		"show", "job_id", jobID,
		ApplicationCommandIdVersion{
			Version: "1169435442328911903",
			ID:      "1169435442328911902",
		})
}

func (c *Client) Show(ctx context.Context, jobID string) (*CommandResult, error) {

	filterFunc := func(msg *discord.Message) (bool, error) {
		if !IsMidjourneyReply(msg) {
			return false, nil
		}

		return MessageJobID(msg) == jobID, nil
	}

	ob := NewCommonMessageObserver(filterFunc)
	c.RegisterMessageObserver(ob)
	defer c.UnregisterMessageObserver(ob)

	payload := NewShowCommandPayload(c.GuildID, c.ChannelID, c.ws.SessionID(), jobID)
	err := c.sendHttpInteraction(ctx, payload)
	if err != nil {
		return nil, err
	}

	msg, err := ob.WaitMsg(time.Minute * 5)
	if err != nil {
		return nil, err
	}

	return MessageCommandResult(msg), nil
}
