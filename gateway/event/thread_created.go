package event

import (
	"github.com/bytedance/sonic"
	"github.com/kiwioneone/gojourney/discord"
	"github.com/kiwioneone/gojourney/rest"
)

type ThreadCreate struct {
	Data *discord.Channel `json:"d"`
}

func NewThreadCreate(rest *rest.Client, data []byte) (*ThreadCreate, error) {
	pk := new(ThreadCreate)

	err := sonic.Unmarshal(data, pk)

	if err != nil {
		return nil, err
	}

	return pk, nil
}
