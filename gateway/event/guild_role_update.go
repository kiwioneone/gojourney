package event

import (
	"github.com/bytedance/sonic"
	"github.com/kiwioneone/gojourney/discord"
	"github.com/kiwioneone/gojourney/rest"
)

// GuildRoleUpdate Is sent when a guild role is updated.
type GuildRoleUpdate struct {
	Data *discord.GuildRoleUpdateEventFields `json:"d"`
}

func NewGuildRoleUpdate(rest *rest.Client, data []byte) (*GuildRoleUpdate, error) {
	pk := new(GuildRoleUpdate)

	err := sonic.Unmarshal(data, pk)

	if err != nil {
		return nil, err
	}

	return pk, nil
}
