package event

import (
	"github.com/bytedance/sonic"
	"github.com/kiwioneone/gojourney/discord"
	"github.com/kiwioneone/gojourney/rest"
)

type GuildEmojisUpdate struct {
	Data *discord.GuildEmojisUpdateEventFields `json:"d"`
}

func NewGuildEmojisUpdate(rest *rest.Client, data []byte) (*GuildEmojisUpdate, error) {
	pk := new(GuildEmojisUpdate)

	err := sonic.Unmarshal(data, pk)

	if err != nil {
		return nil, err
	}

	return pk, nil
}
