package event

import (
	"github.com/bytedance/sonic"
	"github.com/kiwioneone/gojourney/discord"
	"github.com/kiwioneone/gojourney/rest"
)

type InteractionCreate struct {
	Data *discord.Interaction `json:"d"`
}

func NewInteractionCreate(rest *rest.Client, data []byte) (*InteractionCreate, error) {
	pk := new(InteractionCreate)

	err := sonic.Unmarshal(data, pk)

	if err != nil {
		return nil, err
	}

	return pk, nil
}
