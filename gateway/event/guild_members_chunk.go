package event

import (
	"github.com/bytedance/sonic"
	"github.com/kiwioneone/gojourney/discord"
	"github.com/kiwioneone/gojourney/rest"
)

// GuildMembersChunk Is sent in response to Guild Request Members. You can use the chunk_index and chunk_count to calculate how many chunks are left for your request.
type GuildMembersChunk struct {
	Data *discord.GuildMembersChunkEventFields `json:"d"`
}

func NewGuildMembersChunk(rest *rest.Client, data []byte) (*GuildMembersChunk, error) {
	pk := new(GuildMembersChunk)

	err := sonic.Unmarshal(data, pk)

	if err != nil {
		return nil, err
	}

	return pk, nil
}
