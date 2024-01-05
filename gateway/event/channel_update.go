package event

import (
	"github.com/bytedance/sonic"
	"github.com/kiwioneone/gojourney/discord"
	"github.com/kiwioneone/gojourney/rest"
)

type ChannelUpdate struct {
	Data *discord.Channel `json:"d"`
}

func NewChannelUpdate(rest *rest.Client, data []byte) (*ChannelUpdate, error) {
	pk := new(ChannelUpdate)

	err := sonic.Unmarshal(data, pk)

	if err != nil {
		return nil, err
	}

	return pk, nil
}
