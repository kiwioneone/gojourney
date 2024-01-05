package event

import (
	"github.com/bytedance/sonic"
	"github.com/kiwioneone/gojourney/discord"
	"github.com/kiwioneone/gojourney/rest"
)

type VoiceServerUpdate struct {
	Data *discord.VoiceServerUpdateEventFields `json:"d"`
}

func NewVoiceServerUpdate(rest *rest.Client, data []byte) (*VoiceServerUpdate, error) {
	pk := new(VoiceServerUpdate)

	err := sonic.Unmarshal(data, pk)

	if err != nil {
		return nil, err
	}

	return pk, nil
}
