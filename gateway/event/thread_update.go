package event

import (
	"github.com/bytedance/sonic"
	"github.com/kiwioneone/gojourney/discord"
	"github.com/kiwioneone/gojourney/rest"
)

type ThreadUpdate struct {
	Data *discord.Channel `json:"d"`
}

func NewThreadUpdate(rest *rest.Client, data []byte) (*ThreadUpdate, error) {
	pk := new(ThreadUpdate)

	err := sonic.Unmarshal(data, pk)

	if err != nil {
		return nil, err
	}

	return pk, nil
}
