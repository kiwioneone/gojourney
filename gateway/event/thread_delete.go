package event

import (
	"github.com/bytedance/sonic"
	"github.com/kiwioneone/gojourney/discord"
	"github.com/kiwioneone/gojourney/rest"
)

type ThreadDelete struct {
	Data *discord.Channel `json:"d"`
}

func NewThreadDelete(rest *rest.Client, data []byte) (*ThreadDelete, error) {
	pk := new(ThreadDelete)

	err := sonic.Unmarshal(data, pk)

	if err != nil {
		return nil, err
	}

	return pk, nil
}
