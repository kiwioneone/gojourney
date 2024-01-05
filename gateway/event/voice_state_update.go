package event

import (
	"github.com/bytedance/sonic"
	"github.com/kiwioneone/gojourney/discord"
	"github.com/kiwioneone/gojourney/rest"
)

type VoiceStateUpdate struct {
	Data *discord.VoiceState `json:"d"`
}

func NewVoiceStateUpdate(rest *rest.Client, data []byte) (*VoiceStateUpdate, error) {
	pk := new(VoiceStateUpdate)

	err := sonic.Unmarshal(data, pk)

	if err != nil {
		return nil, err
	}

	return pk, nil
}
