package event

import (
	"github.com/bytedance/sonic"
	"github.com/kiwioneone/gojourney/discord"
	"github.com/kiwioneone/gojourney/rest"
)

type GuildBanRemove struct {
	Data struct {
		GuildId string        `json:"guild_id"`
		User    *discord.User `json:"user"`
	} `json:"d"`
}

func NewGuildBanRemove(rest *rest.Client, data []byte) (*GuildBanRemove, error) {
	pk := new(GuildBanRemove)

	err := sonic.Unmarshal(data, pk)

	if err != nil {
		return nil, err
	}

	return pk, nil
}
