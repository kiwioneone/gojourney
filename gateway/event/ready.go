package event

import (
	"github.com/bytedance/sonic"
	"github.com/kiwioneone/gojourney/discord"
)

type Ready struct {
	Data struct {
		Version   int              `json:"v"`
		User      *discord.User    `json:"user"`
		Guilds    []*discord.Guild `json:"guilds"`
		SessionID string           `json:"session_id"`
		Shard     []int            `json:"shard,omitempty"`
	} `json:"d"`
}

func NewReady(data []byte) (*Ready, error) {
	pk := new(Ready)

	err := sonic.Unmarshal(data, pk)

	if err != nil {
		return nil, err
	}

	return pk, nil
}
