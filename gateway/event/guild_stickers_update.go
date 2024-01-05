package event

import (
	"github.com/bytedance/sonic"
	"github.com/kiwioneone/gojourney/discord"
	"github.com/kiwioneone/gojourney/rest"
)

type GuildStickersUpdate struct {
	Data *discord.GuildStickersUpdateEventFields `json:"d"`
}

func NewGuildStickersUpdate(rest *rest.Client, data []byte) (*GuildStickersUpdate, error) {
	pk := new(GuildStickersUpdate)

	err := sonic.Unmarshal(data, pk)

	if err != nil {
		return nil, err
	}

	return pk, nil
}
