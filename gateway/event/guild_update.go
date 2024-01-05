package event

import (
	"github.com/bytedance/sonic"
	"github.com/kiwioneone/gojourney/discord"
	"github.com/kiwioneone/gojourney/rest"
)

type GuildUpdate struct {
	Data *discord.Guild `json:"d"`
}

func NewGuildUpdate(rest *rest.Client, data []byte) (*GuildUpdate, error) {
	pk := new(GuildUpdate)

	err := sonic.Unmarshal(data, pk)

	if err != nil {
		return nil, err
	}

	return pk, nil
}
