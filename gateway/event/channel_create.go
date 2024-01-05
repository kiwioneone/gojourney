package event

import (
	"github.com/bytedance/sonic"
	"github.com/kiwioneone/gojourney/discord"
	"github.com/kiwioneone/gojourney/rest"
)

type ChannelCreate struct {
	Data *discord.Channel `json:"d"`
}

func NewChannelCreate(rest *rest.Client, data []byte) (*ChannelCreate, error) {
	pk := new(ChannelCreate)

	err := sonic.Unmarshal(data, pk)

	if err != nil {
		return nil, err
	}

	return pk, nil
}
