package gojourney

import (
	"context"
	"fmt"
	"strings"

	"github.com/kiwioneone/gojourney/discord"
)

// AppealBlockAction for 'Action need to continue'
func (c *Client) AppealBlockAction(msg *discord.Message) {
	if msg.Flags != discord.MessageFlagEphemeral {
		return
	}

	if len(msg.Components) == 0 {
		return
	}

	actions, ok := msg.Components[0].(*discord.ActionRows)
	if !ok {
		return
	}

	if len(actions.Components) == 0 {
		return
	}

	button, ok := actions.Components[0].(*discord.Button)
	if !ok {
		return
	}

	if !strings.HasPrefix(button.CustomId, "MJ::Prompts::") {
		return
	}

	payload := NewMessageComponentPayload(c.GuildID, c.ChannelID, c.ws.SessionID(), MessageID(msg), discord.ComponentTypeButton, button.CustomId)
	payload.MessageFlags = int(discord.MessageFlagEphemeral)

	ctx := context.Background()
	err := c.sendHttpInteraction(ctx, payload)
	fmt.Printf("appeal %s, error %#v\n", button.CustomId, err)
}
