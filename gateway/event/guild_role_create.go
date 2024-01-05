package event

import (
	"github.com/bytedance/sonic"
	"github.com/kiwioneone/gojourney/discord"
	"github.com/kiwioneone/gojourney/rest"
)

// GuildRoleCreate Is sent when a guild role is created.
type GuildRoleCreate struct {
	Data *discord.GuildRoleCreateEventFields `json:"d"`
}

func NewGuildRoleCreate(rest *rest.Client, data []byte) (*GuildRoleCreate, error) {
	pk := new(GuildRoleCreate)

	err := sonic.Unmarshal(data, pk)

	if err != nil {
		return nil, err
	}

	return pk, nil
}
