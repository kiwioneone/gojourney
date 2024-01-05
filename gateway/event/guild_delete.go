package event

import (
	"github.com/bytedance/sonic"
	"github.com/kiwioneone/gojourney/discord"
	"github.com/kiwioneone/gojourney/rest"
)

type GuildDelete struct {
	Data *discord.Guild `json:"d"`
}

func NewGuildDelete(rest *rest.Client, data []byte) (*GuildDelete, error) {
	pk := new(GuildDelete)

	err := sonic.Unmarshal(data, pk)

	if err != nil {
		return nil, err
	}

	return pk, nil
}
