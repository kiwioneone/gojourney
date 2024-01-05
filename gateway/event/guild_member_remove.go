package event

import (
	"github.com/bytedance/sonic"
	"github.com/kiwioneone/gojourney/discord"
	"github.com/kiwioneone/gojourney/rest"
)

type GuildMemberRemove struct {
	Data *discord.GuildMemberRemoveEventFields `json:"d"`
}

func NewGuildMemberRemove(rest *rest.Client, data []byte) (*GuildMemberRemove, error) {
	pk := new(GuildMemberRemove)

	err := sonic.Unmarshal(data, pk)

	if err != nil {
		return nil, err
	}

	return pk, nil
}
