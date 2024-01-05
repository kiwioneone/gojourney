package gojourney

import (
	"context"
	"fmt"
	"time"

	"github.com/kiwioneone/gojourney/discord"
)

func (c *Client) ZoomOut2x(ctx context.Context, messageID string, jobID string) (*CommandResult, error) {
	customID := fmt.Sprintf("MJ::Outpaint::50::1::%s::SOLO", jobID)
	return c.runJob(ctx, messageID, customID)
}

func (c *Client) ZoomOut15x(ctx context.Context, messageID string, jobID string) (*CommandResult, error) {
	customID := fmt.Sprintf("MJ::Outpaint::75::1::%s::SOLO", jobID)
	return c.runJob(ctx, messageID, customID)
}

func (c *Client) CustomZoom(ctx context.Context, messageID string, jobID string, prompt string, factor float64) (result *CommandResult, err error) {
	if factor <= 1 || factor > 2 {
		return nil, ErrorFactorInvalid
	}

	// first stage

	customID := fmt.Sprintf("MJ::CustomZoom::%s", jobID)

	nextID, err := c.modalSubmit(ctx, messageID, customID)
	if err != nil {
		return nil, err
	}

	// second stage
	modalCustomID := fmt.Sprintf("MJ::OutpaintCustomZoomModal::%s", jobID)
	components := []discord.MessageComponent{
		discord.ActionRows{
			Type: discord.ComponentTypeActionRow,
			Components: []discord.MessageComponent{
				discord.TextInput{
					Type:     discord.ComponentTypeTextInput,
					CustomId: "MJ::OutpaintCustomZoomModal::prompt",
					Value:    fmt.Sprintf("%s --zoom %f", prompt, factor),
				},
			},
		},
	}

	//filterFunc := func(msg *discord.Message) bool {
	//	if !IsMidjourneyReply(msg) {
	//		return false
	//	}
	//
	//	if ReferenceMessageID(msg) != messageID {
	//		return false
	//	}
	//
	//	if strings.Contains(msg.Content, "Zoom Out") {
	//		return true
	//	}
	//
	//	return false
	//}
	//
	//ob := NewCommonMessageObserver(filterFunc)

	payload := NewModalSubmitPayload(c.GuildID, c.ChannelID, c.ws.SessionID(), nextID, modalCustomID, components)

	messageFinder := &MessageFinder{Nonce: payload.Nonce}
	ob := NewCommonMessageObserver(messageFinder.FilterMessage)
	c.RegisterMessageObserver(ob)
	defer c.UnregisterMessageObserver(ob)

	err = c.sendHttpInteraction(ctx, payload)
	if err != nil {
		return nil, err
	}

	zoomMsg, err := ob.WaitMsg(time.Minute * 5)
	if err != nil {
		return nil, err
	}

	return MessageCommandResult(zoomMsg), nil
}
