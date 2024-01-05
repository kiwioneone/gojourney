package event

import (
	"github.com/bytedance/sonic"
	"github.com/kiwioneone/gojourney/rest"
)

type MessageDeleteData struct {
	Id        string `json:"id"`
	ChannelId string `json:"channel_id"`
	GuildId   string `json:"guild_id"`
}

type MessageDelete struct {
	Data *MessageDeleteData `json:"d"`
}

func NewMessageDelete(_ *rest.Client, data []byte) (*MessageDelete, error) {
	pk := new(MessageDelete)

	err := sonic.Unmarshal(data, pk)

	if err != nil {
		return nil, err
	}

	return pk, nil
}
