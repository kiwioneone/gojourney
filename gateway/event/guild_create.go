package event

import (
	"github.com/bytedance/sonic"
	"github.com/kiwioneone/gojourney/discord"
	"github.com/kiwioneone/gojourney/rest"
)

type GuildCreate struct {
	Data *discord.Guild `json:"d"`
}

func NewGuildCreate(rest *rest.Client, data []byte) (*GuildCreate, error) {
	pk := new(GuildCreate)

	err := sonic.Unmarshal(data, pk)

	if err != nil {
		return nil, err
	}

	return pk, nil
}
