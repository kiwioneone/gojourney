package event

import (
	"github.com/bytedance/sonic"
	"github.com/kiwioneone/gojourney/discord"
	"github.com/kiwioneone/gojourney/rest"
)

// GuildRoleDelete Is sent when a guild role is deleted.
type GuildRoleDelete struct {
	Data *discord.GuildRoleDeleteEventFields `json:"d"`
}

func NewGuildRoleDelete(rest *rest.Client, data []byte) (*GuildRoleDelete, error) {
	pk := new(GuildRoleDelete)

	err := sonic.Unmarshal(data, pk)

	if err != nil {
		return nil, err
	}

	return pk, nil
}
