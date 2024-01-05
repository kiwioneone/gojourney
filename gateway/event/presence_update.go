package event

import (
	"github.com/bytedance/sonic"
	"github.com/kiwioneone/gojourney/discord"
	"github.com/kiwioneone/gojourney/rest"
)

type PresenceUpdate struct {
	Data struct {
		User         *discord.User       `json:"user"`
		GuildId      string              `json:"guild_id"`
		Status       string              `json:"status"`
		Activities   []*discord.Activity `json:"activities"`
		ClientStatus *ClientStatus       `json:"client_status"`
	} `json:"d"`
}

type ClientStatus struct {
	Desktop string `json:"deskop,omitempty"` // windows, linux, mac
	Mobile  string `json:"mobile,omitempty"` // ios, android
	Web     string `json:"web,omitempty"`    // browser, bot_account
}

func NewPresenceUpdate(_ *rest.Client, data []byte) (*PresenceUpdate, error) {
	pk := new(PresenceUpdate)

	err := sonic.Unmarshal(data, pk)

	if err != nil {
		return nil, err
	}

	return pk, nil
}
