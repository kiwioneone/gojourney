package gojourney

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/kiwioneone/gojourney/discord"
)

type modalSubmitObserver struct {
	nonce string
	ch    chan string
}

type ModalSubmitInteraction struct {
	Nonce string `json:"nonce,omitempty"`
	ID    string `json:"id,omitempty"`
}

func (o *modalSubmitObserver) Observe(interaction *discord.Interaction) {
	if interaction == nil {
		return
	}

	var modalSubmitInteraction ModalSubmitInteraction
	err := json.Unmarshal(interaction.RawData, &modalSubmitInteraction)

	if err == nil && modalSubmitInteraction.Nonce == o.nonce {
		o.ch <- modalSubmitInteraction.ID
	}
}

func (o *modalSubmitObserver) Wait(timeout time.Duration) (string, error) {
	defer close(o.ch)

	select {
	case result := <-o.ch:
		return result, nil
	case <-time.After(timeout):
		return "", errors.New("timeout")
	}
}

func (c *Client) modalSubmit(ctx context.Context, messageID string, customID string) (nextID string, err error) {

	payload := NewMessageComponentPayload(c.GuildID, c.ChannelID, c.ws.SessionID(), messageID, discord.ComponentTypeButton, customID)
	nonce := payload.Nonce

	ob := &modalSubmitObserver{nonce: nonce, ch: make(chan string)}
	c.RegisterInteractionObserver(ob)
	defer c.UnregisterInteractionObserver(ob)

	err = c.sendHttpInteraction(ctx, payload)
	if err != nil {
		return "", err
	}

	nextID, err = ob.Wait(time.Minute * 5)
	return nextID, err
}
